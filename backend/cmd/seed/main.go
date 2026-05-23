// Command seed populates the database with realistic data for manual testing.
//
//	go run ./cmd/seed                       # apply curated dataset
//	go run ./cmd/seed -reset                # wipe everything first
//	go run ./cmd/seed -reset -scale 1000    # plus 1000 procedural users + proportional content
//	go run ./cmd/seed -file path/to.json    # use a different curated dataset
//
// Curated data lives in cmd/seed/data/seed.json — references are by nickname /
// group title so you can edit it without thinking about IDs. Procedural bulk
// data is generated synthetically and inserted with batched transactions, so
// SCALE=5000 finishes in seconds and SCALE=50000 finishes in a couple minutes.
//
// All seeded users share the password "password123".

package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"social-network/internal/config"
	"social-network/internal/database"
	"social-network/internal/domain"
	"social-network/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

const seedPassword = "password123"

func main() {
	reset := flag.Bool("reset", false, "wipe all tables before seeding")
	scale := flag.Int("scale", 0, "generate this many extra procedural users plus proportional content")
	file := flag.String("file", "cmd/seed/data/seed.json", "path to curated seed JSON")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.NewDatabase(cfg.Database.Path)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	if *reset {
		fmt.Println("Resetting database...")
		if err := resetTables(db); err != nil {
			log.Fatalf("reset tables: %v", err)
		}
	}

	rand.Seed(time.Now().UnixNano())
	repos := repository.NewRepositories(db)

	fmt.Println("Generating seed images...")
	avatars, postMedia, err := generateSeedImages(cfg.Upload.UploadPath)
	if err != nil {
		log.Fatalf("generate images: %v", err)
	}
	fmt.Printf("  %d avatar PNGs, %d post media PNGs in %s\n", len(avatars), len(postMedia), cfg.Upload.UploadPath)

	curated, err := loadSeed(*file)
	if err != nil {
		log.Fatalf("load seed file %s: %v", *file, err)
	}

	fmt.Println("Applying curated seed...")
	nickToID, err := applyCurated(db, repos, curated, avatars, postMedia)
	if err != nil {
		log.Fatalf("apply curated: %v", err)
	}
	fmt.Printf("  %d users, %d posts, %d direct chats, %d groups, %d notifications\n",
		len(curated.Users), len(curated.Posts), len(curated.DirectChats), len(curated.Groups), len(curated.Notifications))

	if *scale > 0 {
		fmt.Printf("Generating procedural bulk: %d users with realistic activity distribution...\n", *scale)
		started := time.Now()
		stats, err := generateBulk(db, *scale, valuesOf(nickToID), avatars, postMedia)
		if err != nil {
			log.Fatalf("generate bulk: %v", err)
		}
		fmt.Printf("  tiers: %d celebrities · %d actives · %d normals · %d lurkers\n",
			stats.celebs, stats.actives, stats.normals, stats.lurkers)
		fmt.Printf("  +%d users, +%d posts, +%d comments, +%d follows, +%d conversations, +%d messages in %s\n",
			stats.users, stats.posts, stats.comments, stats.follows, stats.convs, stats.messages,
			time.Since(started).Round(time.Millisecond))
	}

	if err := finalizeSeed(db); err != nil {
		log.Fatalf("finalize seed: %v", err)
	}

	fmt.Println()
	fmt.Println("Done. Sample logins:")
	fmt.Println("  email: alice@example.com   password: password123")
	fmt.Println("  email: bob@example.com     password: password123")
	fmt.Println("  nickname: cgaldan          password: <any> (dev bypass)")
	fmt.Println("  nickname: testuser         password: <any> (dev bypass)")
}

// ---------------- JSON schema ----------------

type seedFile struct {
	Users         []userSpec         `json:"users"`
	Follows       []followSpec       `json:"follows"`
	Posts         []postSpec         `json:"posts"`
	DirectChats   []directChatSpec   `json:"direct_chats"`
	Groups        []groupSpec        `json:"groups"`
	Notifications []notificationSpec `json:"notifications"`
}

type userSpec struct {
	Nickname    string `json:"nickname"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	AboutMe     string `json:"about_me"`
	IsPublic    bool   `json:"is_public"`
	DateOfBirth string `json:"date_of_birth"`
	AvatarPath  string `json:"avatar_path"`
}

type followSpec struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Status string `json:"status"`
}

type commentSpec struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type postSpec struct {
	Author   string        `json:"author"`
	Title    string        `json:"title"`
	Content  string        `json:"content"`
	Category string        `json:"category"`
	Privacy  string        `json:"privacy"`
	Audience []string      `json:"audience"`
	Comments []commentSpec `json:"comments"`
}

type directChatSpec struct {
	Between  []string      `json:"between"`
	Messages []messageSpec `json:"messages"`
}

type messageSpec struct {
	From string `json:"from"`
	Text string `json:"text"`
}

type groupSpec struct {
	Creator      string          `json:"creator"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Members      []string        `json:"members"`
	Invited      []string        `json:"invited"`
	JoinRequests []string        `json:"join_requests"`
	Posts        []groupPostSpec `json:"posts"`
	Events       []eventSpec     `json:"events"`
}

type groupPostSpec struct {
	Author   string `json:"author"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type eventSpec struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	StartsInDays int        `json:"starts_in_days"`
	Creator      string     `json:"creator"`
	RSVPs        []rsvpSpec `json:"rsvps"`
}

type rsvpSpec struct {
	User     string `json:"user"`
	Response string `json:"response"`
}

type notificationSpec struct {
	Recipient string `json:"recipient"`
	Actor     string `json:"actor"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Body      string `json:"body"`
}

func loadSeed(path string) (*seedFile, error) {
	if !filepath.IsAbs(path) {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// allow running from repo root or backend dir
			alt := filepath.Join("backend", path)
			if _, err2 := os.Stat(alt); err2 == nil {
				path = alt
			}
		}
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var sf seedFile
	if err := json.Unmarshal(raw, &sf); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}
	return &sf, nil
}

// ---------------- Reset ----------------

// Listed child-to-parent so FK cascades stay happy.
var resetOrder = []string{
	"notifications",
	"group_event_rsvps",
	"group_events",
	"group_invitations",
	"group_join_requests",
	"group_members",
	"groups",
	"messages",
	"conversation_participants",
	"conversations",
	"post_audiences",
	"comments",
	"posts",
	"follows",
	"sessions",
	"users",
}

// finalizeSeed cleans up two cosmetic quirks so the DB inspector sees what the
// API sees:
//   - is_online defaults to TRUE in the schema; reset every seeded user to
//     offline since no one is actually connected.
//   - following_count / followers_count are stored on users but the app
//     ignores them and recomputes via subqueries. Write the real counts so
//     raw DB queries match the live UI at seed time.
func finalizeSeed(db *sql.DB) error {
	if _, err := db.Exec("UPDATE users SET is_online = 0"); err != nil {
		return fmt.Errorf("reset is_online: %w", err)
	}
	if _, err := db.Exec(`
		UPDATE users SET
			following_count = (SELECT COUNT(*) FROM follows WHERE follower_id  = users.id AND status = 'accepted'),
			followers_count = (SELECT COUNT(*) FROM follows WHERE following_id = users.id AND status = 'accepted')`); err != nil {
		return fmt.Errorf("snapshot follow counts: %w", err)
	}
	return nil
}

func resetTables(db *sql.DB) error {
	for _, table := range resetOrder {
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			return fmt.Errorf("delete from %s: %w", table, err)
		}
	}
	_, _ = db.Exec("DELETE FROM sqlite_sequence")
	return nil
}

// ---------------- Curated seed (uses repositories for correctness) ----------------

func applyCurated(db *sql.DB, repos *repository.Repositories, sf *seedFile, avatars, postMedia []string) (map[string]int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	hashStr := string(hash)

	nickToID := make(map[string]int, len(sf.Users))
	for i, u := range sf.Users {
		dob, err := parseDate(u.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("user %s: %w", u.Nickname, err)
		}
		avatar := u.AvatarPath
		if avatar == "" && len(avatars) > 0 {
			avatar = avatars[i%len(avatars)]
		}
		id, err := repos.User.CreateUser(
			u.Email, hashStr, u.FirstName, u.LastName, dob,
			u.Nickname, u.Gender, avatar, u.AboutMe, u.IsPublic,
		)
		if err != nil {
			return nil, fmt.Errorf("create user %s: %w", u.Nickname, err)
		}
		nickToID[u.Nickname] = int(id)
	}

	mustID := func(nick string) (int, error) {
		id, ok := nickToID[nick]
		if !ok {
			return 0, fmt.Errorf("unknown nickname %q", nick)
		}
		return id, nil
	}

	for _, f := range sf.Follows {
		from, err := mustID(f.From)
		if err != nil {
			return nil, err
		}
		to, err := mustID(f.To)
		if err != nil {
			return nil, err
		}
		if _, err := repos.Follow.CreateFollow(from, to, f.Status); err != nil {
			return nil, fmt.Errorf("follow %s->%s: %w", f.From, f.To, err)
		}
	}

	for i, p := range sf.Posts {
		author, err := mustID(p.Author)
		if err != nil {
			return nil, err
		}
		media := ""
		// give roughly 1 in 3 curated posts an image, deterministically
		if len(postMedia) > 0 && i%3 == 0 {
			media = postMedia[i%len(postMedia)]
		}
		postID, err := repos.Post.CreatePost(author, p.Title, p.Content, p.Category, p.Privacy, media, 0)
		if err != nil {
			return nil, fmt.Errorf("create post %q: %w", p.Title, err)
		}
		if p.Privacy == "almost_private" && len(p.Audience) > 0 {
			audience := make([]int, 0, len(p.Audience))
			for _, nick := range p.Audience {
				id, err := mustID(nick)
				if err != nil {
					return nil, err
				}
				audience = append(audience, id)
			}
			if err := repos.Post.SetPostAudience(int(postID), audience); err != nil {
				return nil, fmt.Errorf("set audience for %q: %w", p.Title, err)
			}
		}
		for _, c := range p.Comments {
			commenter, err := mustID(c.Author)
			if err != nil {
				return nil, err
			}
			if _, err := repos.Comment.CreateComment(commenter, int(postID), c.Content, ""); err != nil {
				return nil, fmt.Errorf("comment on %q: %w", p.Title, err)
			}
		}
	}

	for _, chat := range sf.DirectChats {
		if len(chat.Between) != 2 {
			return nil, fmt.Errorf("direct chat must have 2 participants, got %d", len(chat.Between))
		}
		a, err := mustID(chat.Between[0])
		if err != nil {
			return nil, err
		}
		b, err := mustID(chat.Between[1])
		if err != nil {
			return nil, err
		}
		conv, err := repos.Conversation.CreateDirectConversation(a, b)
		if err != nil {
			return nil, fmt.Errorf("direct conv %s<->%s: %w", chat.Between[0], chat.Between[1], err)
		}
		for _, m := range chat.Messages {
			sender, err := mustID(m.From)
			if err != nil {
				return nil, err
			}
			if _, err := repos.Message.CreateMessage(&domain.Message{
				ConversationID: conv.ID,
				SenderID:       sender,
				Content:        m.Text,
			}); err != nil {
				return nil, fmt.Errorf("message: %w", err)
			}
		}
	}

	now := time.Now()
	for _, g := range sf.Groups {
		creatorID, err := mustID(g.Creator)
		if err != nil {
			return nil, err
		}
		memberIDs := make([]int, 0, len(g.Members))
		for _, nick := range g.Members {
			id, err := mustID(nick)
			if err != nil {
				return nil, err
			}
			memberIDs = append(memberIDs, id)
		}
		conv, err := repos.Conversation.CreateGroupConversation(g.Title, memberIDs...)
		if err != nil {
			return nil, fmt.Errorf("group conv %q: %w", g.Title, err)
		}
		groupID, err := repos.Group.CreateGroup(&domain.Group{
			CreatorID:      creatorID,
			Title:          g.Title,
			Description:    g.Description,
			ConversationID: conv.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("create group %q: %w", g.Title, err)
		}
		for _, id := range memberIDs {
			role := "member"
			if id == creatorID {
				role = "admin"
			}
			if err := repos.Group.AddMember(int(groupID), id, role); err != nil {
				return nil, fmt.Errorf("add member to %q: %w", g.Title, err)
			}
		}
		for _, nick := range g.Invited {
			id, err := mustID(nick)
			if err != nil {
				return nil, err
			}
			if _, err := repos.Group.CreateGroupInvitation(int(groupID), creatorID, id); err != nil {
				return nil, fmt.Errorf("invite to %q: %w", g.Title, err)
			}
		}
		for _, nick := range g.JoinRequests {
			id, err := mustID(nick)
			if err != nil {
				return nil, err
			}
			if _, err := repos.Group.CreateGroupJoinRequest(int(groupID), id); err != nil {
				return nil, fmt.Errorf("join request %q: %w", g.Title, err)
			}
		}
		for _, p := range g.Posts {
			author, err := mustID(p.Author)
			if err != nil {
				return nil, err
			}
			if _, err := repos.Post.CreatePost(author, p.Title, p.Content, p.Category, "public", "", int(groupID)); err != nil {
				return nil, fmt.Errorf("group post %q: %w", p.Title, err)
			}
		}
		for _, ev := range g.Events {
			evCreator := creatorID
			if ev.Creator != "" {
				id, err := mustID(ev.Creator)
				if err != nil {
					return nil, err
				}
				evCreator = id
			}
			eventID, err := repos.Group.CreateGroupEvent(&domain.GroupEvent{
				GroupID:     int(groupID),
				CreatorID:   evCreator,
				Title:       ev.Title,
				Description: ev.Description,
				StartsAt:    now.Add(time.Duration(ev.StartsInDays) * 24 * time.Hour),
			})
			if err != nil {
				return nil, fmt.Errorf("create event %q: %w", ev.Title, err)
			}
			for _, r := range ev.RSVPs {
				uid, err := mustID(r.User)
				if err != nil {
					return nil, err
				}
				if err := repos.Group.SetGroupEventRSVP(int(eventID), uid, r.Response); err != nil {
					return nil, fmt.Errorf("rsvp: %w", err)
				}
			}
		}
	}

	for _, n := range sf.Notifications {
		recipient, err := mustID(n.Recipient)
		if err != nil {
			return nil, err
		}
		var actorPtr *int
		if n.Actor != "" {
			id, err := mustID(n.Actor)
			if err != nil {
				return nil, err
			}
			actorPtr = &id
		}
		if _, err := repos.Notification.CreateNotification(&domain.Notification{
			RecipientID: recipient,
			ActorID:     actorPtr,
			Type:        n.Type,
			Title:       n.Title,
			Body:        n.Body,
		}); err != nil {
			return nil, fmt.Errorf("notification: %w", err)
		}
	}

	return nickToID, nil
}

func parseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Date(1995, 6, 1, 0, 0, 0, 0, time.UTC), nil
	}
	return time.Parse("2006-01-02", s)
}

func valuesOf(m map[string]int) []int {
	out := make([]int, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

// ---------------- Procedural bulk with tiered distribution ----------------

type bulkStats struct {
	users, posts, comments, follows, convs, messages int
	celebs, actives, normals, lurkers                int
}

// tierProfile defines how a single user in a tier behaves: how many posts they
// write, who they follow, etc. All ranges are inclusive of lo, exclusive of hi.
type tierProfile struct {
	name                          string
	postsLo, postsHi              int
	outFollowsLo, outFollowsHi    int
	convsLo, convsHi              int
	msgsPerConvLo, msgsPerConvHi  int
	commentsLo, commentsHi        int // outbound comments per user
	avatarChance, postMediaChance float64
}

var (
	lurkerProfile = tierProfile{
		name:    "lurker",
		postsLo: 0, postsHi: 2,
		outFollowsLo: 0, outFollowsHi: 2,
		convsLo: 0, convsHi: 0,
		msgsPerConvLo: 0, msgsPerConvHi: 0,
		commentsLo: 0, commentsHi: 2,
		avatarChance: 0.20, postMediaChance: 0.02,
	}
	normalProfile = tierProfile{
		name:    "normal",
		postsLo: 1, postsHi: 6,
		outFollowsLo: 3, outFollowsHi: 12,
		convsLo: 0, convsHi: 3,
		msgsPerConvLo: 2, msgsPerConvHi: 6,
		commentsLo: 2, commentsHi: 8,
		avatarChance: 0.70, postMediaChance: 0.15,
	}
	activeProfile = tierProfile{
		name:    "active",
		postsLo: 5, postsHi: 30,
		outFollowsLo: 20, outFollowsHi: 80,
		convsLo: 3, convsHi: 8,
		msgsPerConvLo: 4, msgsPerConvHi: 12,
		commentsLo: 10, commentsHi: 40,
		avatarChance: 0.95, postMediaChance: 0.35,
	}
	celebProfile = tierProfile{
		name:    "celebrity",
		postsLo: 30, postsHi: 200,
		outFollowsLo: 20, outFollowsHi: 80,
		convsLo: 5, convsHi: 15,
		msgsPerConvLo: 6, msgsPerConvHi: 18,
		commentsLo: 20, commentsHi: 80,
		avatarChance: 1.00, postMediaChance: 0.55,
	}
)

// generateBulk inserts scale synthetic users with a realistic activity mix
// (60% lurkers, 30% normals, ~9% actives, ~1% celebrities) plus matching
// posts, comments, follows, conversations and messages. Celebrities also act
// as follow magnets, attracting follows from a large random share of the
// population. All inserts happen inside one transaction for speed.
//
// existingUserIDs lets follows, comments and chats link to curated users too.
func generateBulk(db *sql.DB, scale int, existingUserIDs []int, avatars, postMedia []string) (bulkStats, error) {
	var stats bulkStats
	if scale <= 0 {
		return stats, nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		return stats, fmt.Errorf("hash: %w", err)
	}
	hashStr := string(hash)

	// --- decide tier counts ---------------------------------------------
	stats.celebs = clamp(scale/500, 1, 50)
	stats.actives = scale * 9 / 100
	stats.normals = scale * 30 / 100
	stats.lurkers = scale - stats.celebs - stats.actives - stats.normals
	if stats.lurkers < 0 {
		stats.lurkers = 0
	}

	// assignment list: a tierProfile per user index, then shuffle so the
	// celebrity IDs aren't clustered at the start of the auto-increment range.
	assign := make([]*tierProfile, 0, scale)
	for i := 0; i < stats.celebs; i++ {
		assign = append(assign, &celebProfile)
	}
	for i := 0; i < stats.actives; i++ {
		assign = append(assign, &activeProfile)
	}
	for i := 0; i < stats.normals; i++ {
		assign = append(assign, &normalProfile)
	}
	for i := 0; i < stats.lurkers; i++ {
		assign = append(assign, &lurkerProfile)
	}
	rand.Shuffle(len(assign), func(i, j int) { assign[i], assign[j] = assign[j], assign[i] })

	tx, err := db.Begin()
	if err != nil {
		return stats, fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback()

	// --- 1. users -------------------------------------------------------
	uStmt, err := tx.Prepare(`
		INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth, nickname, gender, avatar_path, about_me, is_public)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return stats, fmt.Errorf("prepare users: %w", err)
	}
	bulkIDs := make([]int, 0, scale)
	tierByUser := make(map[int]*tierProfile, scale)
	celebIDs := make([]int, 0, stats.celebs)
	dob := time.Date(1995, 6, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < scale; i++ {
		profile := assign[i]
		nick := fmt.Sprintf("user_%05d", i+1)
		email := nick + "@seed.local"
		about := fmt.Sprintf("Seeded %s user.", profile.name)
		avatar := ""
		if len(avatars) > 0 && rand.Float64() < profile.avatarChance {
			avatar = avatars[rand.Intn(len(avatars))]
		}
		isPublic := profile != &celebProfile || true // celebrities are always public
		res, err := uStmt.Exec(email, hashStr, randomFirst(i), randomLast(i), dob, nick, randomGender(i), avatar, about, isPublic)
		if err != nil {
			uStmt.Close()
			return stats, fmt.Errorf("insert user %d: %w", i, err)
		}
		id, _ := res.LastInsertId()
		uid := int(id)
		bulkIDs = append(bulkIDs, uid)
		tierByUser[uid] = profile
		if profile == &celebProfile {
			celebIDs = append(celebIDs, uid)
		}
	}
	uStmt.Close()
	stats.users = scale

	allUserIDs := append(append([]int{}, existingUserIDs...), bulkIDs...)

	// --- 2. posts -------------------------------------------------------
	pStmt, err := tx.Prepare(`
		INSERT INTO posts (user_id, title, content, category, privacy_level, media_url)
		VALUES (?, ?, ?, ?, 'public', ?)`)
	if err != nil {
		return stats, fmt.Errorf("prepare posts: %w", err)
	}
	postIDs := make([]int, 0, scale*5)
	categories := []string{"general", "tech", "sports", "books", "food", "music", "diy", "garden", "photography", "city"}
	for _, uid := range bulkIDs {
		profile := tierByUser[uid]
		n := randIn(profile.postsLo, profile.postsHi)
		for j := 0; j < n; j++ {
			cat := categories[rand.Intn(len(categories))]
			title := fmt.Sprintf("%s post #%d", cat, j+1)
			media := ""
			if len(postMedia) > 0 && rand.Float64() < profile.postMediaChance {
				media = postMedia[rand.Intn(len(postMedia))]
			}
			res, err := pStmt.Exec(uid, title, loremShort(), cat, media)
			if err != nil {
				pStmt.Close()
				return stats, fmt.Errorf("insert post: %w", err)
			}
			pid, _ := res.LastInsertId()
			postIDs = append(postIDs, int(pid))
		}
	}
	pStmt.Close()
	stats.posts = len(postIDs)

	// --- 3. comments ----------------------------------------------------
	if len(postIDs) > 0 {
		cStmt, err := tx.Prepare(`INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`)
		if err != nil {
			return stats, fmt.Errorf("prepare comments: %w", err)
		}
		for _, uid := range bulkIDs {
			profile := tierByUser[uid]
			n := randIn(profile.commentsLo, profile.commentsHi)
			for k := 0; k < n; k++ {
				postID := postIDs[rand.Intn(len(postIDs))]
				if _, err := cStmt.Exec(postID, uid, loremTiny()); err != nil {
					cStmt.Close()
					return stats, fmt.Errorf("insert comment: %w", err)
				}
				stats.comments++
			}
		}
		cStmt.Close()
	}

	// --- 4. outbound follows -------------------------------------------
	fStmt, err := tx.Prepare(`INSERT OR IGNORE INTO follows (follower_id, following_id, status) VALUES (?, ?, 'accepted')`)
	if err != nil {
		return stats, fmt.Errorf("prepare follows: %w", err)
	}
	for _, uid := range bulkIDs {
		profile := tierByUser[uid]
		n := randIn(profile.outFollowsLo, profile.outFollowsHi)
		for j := 0; j < n; j++ {
			target := allUserIDs[rand.Intn(len(allUserIDs))]
			if target == uid {
				continue
			}
			res, err := fStmt.Exec(uid, target)
			if err != nil {
				fStmt.Close()
				return stats, fmt.Errorf("insert follow: %w", err)
			}
			if n, _ := res.RowsAffected(); n > 0 {
				stats.follows++
			}
		}
	}

	// --- 4b. celebrity inbound follows (the magnet effect) -------------
	// each celebrity is followed by random N/20..N/5 of all users
	if len(celebIDs) > 0 {
		minIn := scale / 20
		maxIn := scale / 5
		if minIn < 10 {
			minIn = 10
		}
		if maxIn < minIn+1 {
			maxIn = minIn + 1
		}
		if maxIn > len(allUserIDs)-1 {
			maxIn = len(allUserIDs) - 1
		}
		for _, celeb := range celebIDs {
			n := randIn(minIn, maxIn)
			for k := 0; k < n; k++ {
				follower := allUserIDs[rand.Intn(len(allUserIDs))]
				if follower == celeb {
					continue
				}
				res, err := fStmt.Exec(follower, celeb)
				if err != nil {
					fStmt.Close()
					return stats, fmt.Errorf("celeb follow: %w", err)
				}
				if n, _ := res.RowsAffected(); n > 0 {
					stats.follows++
				}
			}
		}
	}
	fStmt.Close()

	// --- 5. conversations + messages -----------------------------------
	conStmt, err := tx.Prepare(`INSERT OR IGNORE INTO conversations (type, pair_key) VALUES ('private', ?)`)
	if err != nil {
		return stats, fmt.Errorf("prepare convs: %w", err)
	}
	partStmt, err := tx.Prepare(`INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?)`)
	if err != nil {
		conStmt.Close()
		return stats, fmt.Errorf("prepare participants: %w", err)
	}
	msgStmt, err := tx.Prepare(`INSERT INTO messages (conversation_id, sender_id, content) VALUES (?, ?, ?)`)
	if err != nil {
		conStmt.Close()
		partStmt.Close()
		return stats, fmt.Errorf("prepare messages: %w", err)
	}
	closeChat := func() {
		conStmt.Close()
		partStmt.Close()
		msgStmt.Close()
	}
	for _, uid := range bulkIDs {
		profile := tierByUser[uid]
		convs := randIn(profile.convsLo, profile.convsHi)
		for c := 0; c < convs; c++ {
			other := allUserIDs[rand.Intn(len(allUserIDs))]
			if other == uid {
				continue
			}
			key := pairKey(uid, other)
			res, err := conStmt.Exec(key)
			if err != nil {
				closeChat()
				return stats, fmt.Errorf("insert conv: %w", err)
			}
			if n, _ := res.RowsAffected(); n == 0 {
				continue
			}
			convID, _ := res.LastInsertId()
			if _, err := partStmt.Exec(convID, uid); err != nil {
				closeChat()
				return stats, fmt.Errorf("participant: %w", err)
			}
			if _, err := partStmt.Exec(convID, other); err != nil {
				closeChat()
				return stats, fmt.Errorf("participant: %w", err)
			}
			stats.convs++
			msgs := randIn(profile.msgsPerConvLo, profile.msgsPerConvHi)
			for m := 0; m < msgs; m++ {
				sender := uid
				if m%2 == 1 {
					sender = other
				}
				if _, err := msgStmt.Exec(convID, sender, loremTiny()); err != nil {
					closeChat()
					return stats, fmt.Errorf("insert msg: %w", err)
				}
				stats.messages++
			}
		}
	}
	closeChat()

	if err := tx.Commit(); err != nil {
		return stats, fmt.Errorf("commit: %w", err)
	}
	return stats, nil
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func randIn(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo)
}

func pairKey(a, b int) string {
	if a > b {
		return fmt.Sprintf("%d-%d", a, b)
	}
	return fmt.Sprintf("%d-%d", b, a)
}

// ---------------- Random helpers ----------------

var (
	firsts  = []string{"Alex", "Sam", "Jordan", "Taylor", "Casey", "Riley", "Morgan", "Quinn", "Avery", "Jamie", "Robin", "Drew", "Skyler", "Reese", "Parker", "Rowan", "Hayden", "Emerson", "Sage", "Phoenix"}
	lasts   = []string{"Smith", "Garcia", "Park", "Nguyen", "Rossi", "Petrov", "Hassan", "Kowalski", "Singh", "Tanaka", "Andersen", "Costa", "Diaz", "Adler", "Reyes", "Vance", "Ito", "Hale", "Brooks", "Quinn"}
	genders = []string{"male", "female", "other"}
	tiny    = []string{"+1", "agreed", "Same.", "Where is this?", "I'll be there.", "Sounds good.", "Saving this.", "Send link?", "Love it.", "Hard pass.", "Maybe next time.", "Will reply later."}
	short   = []string{
		"Quick thought, mostly thinking out loud here.",
		"Wrote this on the train. Will follow up later this week.",
		"Trying a thing. Will report back if it works.",
		"Anyone got a recommendation for this?",
		"Long week. Heading to the coast tomorrow.",
		"Posting from the new laptop, finally set up.",
		"Took the long way home and it was worth it.",
		"Tiny win: finally fixed the flaky test.",
		"Reading something interesting, will share next week.",
		"Open question: how do you handle this kind of edge case?",
	}
)

func randomFirst(seed int) string { return firsts[seed%len(firsts)] }
func randomLast(seed int) string  { return lasts[(seed*7)%len(lasts)] }
func randomGender(seed int) string {
	return genders[seed%len(genders)]
}

// randBool returns true with the given probability, seeded for variety.
func randBool(seed int, p float64) bool {
	r := float64((seed*2654435761)%1000) / 1000.0
	return r < p
}

func loremShort() string {
	return short[rand.Intn(len(short))]
}

func loremTiny() string {
	return tiny[rand.Intn(len(tiny))]
}
