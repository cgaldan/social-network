cd backend
go mod tidy
go run cmd/server/main.go

new terminal:
cd frontend
npm install
node ./node_modules/next/dist/bin/next dev
open localhost:3000