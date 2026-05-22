# Frontend Architecture

## Overview

The frontend has been refactored from a monolithic single-file architecture into a modern, modular structure using ES6 modules. This improves maintainability, testability, and code reusability.

## Directory Structure

```
frontend/
├── js/
│   ├── config.js              # Application configuration
│   ├── main.js                # Application entry point
│   ├── api/
│   │   └── client.js          # API client abstraction layer
│   ├── modules/               # Feature modules
│   │   ├── auth.js            # Authentication logic
│   │   ├── posts.js           # Posts management
│   │   ├── messages.js        # Messaging/chat logic
│   │   ├── websocket.js       # WebSocket connection
│   │   └── index.js           # Module barrel exports
│   ├── state/
│   │   └── store.js           # Centralized state management
│   ├── ui/
│   │   ├── ui.js              # UI utilities and rendering
│   │   └── events.js          # Event listener setup
│   └── utils/
│       └── helpers.js         # Helper functions and utilities
├── styles.css                 # Application styles
├── index.html                 # HTML entry point
└── app.js                     # DEPRECATED - Old monolithic file (kept for reference)
```

## Key Improvements

### 1. **Separation of Concerns**
- **API Layer**: All HTTP requests centralized in `api/client.js`
- **State Management**: Single source of truth in `state/store.js`
- **UI Module**: DOM manipulation separated from business logic
- **Feature Modules**: Each feature (auth, posts, messages) in own file

### 2. **Centralized State Management**
The `Store` class provides:
- Single state object for entire application
- `setState()` for updates with change notifications
- `subscribe()` for reactive updates
- Clear getter/setter methods

### 3. **API Client Abstraction**
`ApiClient` class:
- Centralized fetch wrapper
- Automatic header management
- Token handling
- Error logging
- Organized method groups by feature

### 4. **Configuration Management**
All constants in one place:
```javascript
export const CONFIG = {
    API_URL: '/api',
    WS_RECONNECT_ATTEMPTS: 5,
    TOAST_DURATION: 3000,
    // ... more configs
};
```

### 5. **Better Error Handling**
- Consistent error display patterns
- User-friendly error messages
- Console logging for debugging
- Proper try/catch in async operations

### 6. **Improved Code Quality**
- JSDoc comments on all functions
- Consistent naming conventions
- ES6 modules for proper encapsulation
- Helper functions to reduce duplication

## Module Responsibilities

### `auth.js`
- Login/register forms
- Session verification
- Logout and session cleanup
- Auth state management

### `posts.js`
- Load and render posts
- Create new posts
- Open post details
- Handle comments
- Category filtering

### `messages.js`
- Load conversations
- Open chat panels
- Load message history (with infinite scroll)
- Send messages
- WebSocket message handling

### `websocket.js`
- WebSocket connection management
- Auto-reconnect logic
- Message routing
- Connection status updates

### `ui.js`
- Error/toast notifications
- UI rendering helpers
- Tab switching
- General DOM utilities

### `store.js`
- Centralized application state
- Reactive state updates
- Subscriber notification system

### `api/client.js`
- HTTP request wrapper
- Automatic token injection
- Organized API methods by feature
- Error handling

## Migration from Old Code

### Old Way (app.js)
```javascript
// Mixed concerns
const state = { currentUser: null, token: null, posts: [] };
async function handleLogin(e) { ... }
async function loadPosts() { ... }
function renderPosts() { ... }
```

### New Way
```javascript
// auth.js
export async function handleLogin(e) { ... }

// posts.js
export async function loadPosts() { ... }
export function renderPosts() { ... }

// store.js
store.set('token', token);
store.setState({ posts: [] });

// api/client.js
api.login(identifier, password);
api.getPosts();
```

## Using the Modules

### Importing
```javascript
// In any module
import store from '../state/store.js';
import api from '../api/client.js';
import { loadPosts } from '../modules/posts.js';
```

### Adding a New Feature

1. Create new module: `js/modules/feature.js`
2. Implement feature functions
3. Export public functions
4. Add event listeners in `ui/events.js`
5. Use `store.js` for state
6. Use `api/client.js` for API calls

Example:
```javascript
// js/modules/feature.js
import api from '../api/client.js';
import store from '../state/store.js';

export async function initFeature() {
    const data = await api.getFeatureData();
    store.set('featureData', data);
}
```

## Performance Considerations

- **Lazy Loading**: Feature modules imported on-demand via dynamic imports
- **Event Delegation**: Single listeners for multiple elements
- **Debouncing/Throttling**: Utilities available for expensive operations
- **State Efficiency**: Only relevant state kept in memory

## Testing

With this modular structure, testing is easier:
- **Unit Tests**: Test individual modules without imports
- **Integration Tests**: Test module interactions through store and API
- **E2E Tests**: Same as before

Example test:
```javascript
import { handleLogin } from '../modules/auth.js';
import api from '../api/client.js';

jest.mock('../api/client.js');

describe('auth', () => {
    it('should set token on login', async () => {
        api.login.mockResolvedValue({ success: true, token: 'abc' });
        // test...
    });
});
```

## Browser Compatibility

Uses ES6 modules which require:
- Chrome 63+
- Firefox 67+
- Safari 11+
- Edge 79+

For older browsers, use a bundler like Webpack or Vite.

## Next Steps

1. Remove old `app.js` after verification
2. Consider using a bundler (Vite/Webpack) for production
3. Add unit tests for each module
4. Add TypeScript for better type safety
5. Consider using a router for multi-page support
