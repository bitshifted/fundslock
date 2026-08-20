# Frontend Agent Guidelines & Runbook

This file contains instructions for working with the **FundsLock frontend** component — a Vue 3 SPA that powers the user-facing interface of the decentralized escrow system.

## Tech Stack

| Category | Tool | Version |
|---|---|---|
| UI Framework | Vue.js | ^3.5.40 (Composition API, `<script setup>`) |
| Build Tool | Vite | ^8.1.5 |
| Routing | vue-router | ^5.2.0 |
| State Management | Pinia | ^4.0.3 |
| CSS Framework | Bootstrap + bootstrap-vue-next | ^5.3.8 / ^0.45.10 |
| Web3 / Wallet | @reown/appkit + ethers + viem | ^1.8.23 / ^6.17.0 / ^2.55.17 |
| Auth Flow | SIWE (Sign-In With Ethereum) | @reown/appkit-siwe + siwe |
| Testing | Vitest + jsdom + @vue/test-utils | ^4.1.10 |
| Linting | ESLint + oxlint | ^10.7.0 / ~1.74.0 |
| Formatting | Prettier | 3.9.5 |

## Project Structure

```
frontend/
├── src/
│   ├── App.vue                     # Root layout (NavBar + RouterView)
│   ├── main.js                     # Application entry point
│   ├── assets/                     # Global styles, logo
│   │   ├── base.css                # CSS variables (light/dark theme)
│   │   ├── main.css                # Global styles
│   │   └── logo.svg
│   ├── auth/                       # Web3 wallet authentication
│   │   ├── appkit.js               # Reown AppKit configuration
│   │   ├── authenticate.js         # Backend auth API calls (nonce, verify, refresh, session)
│   │   └── siwe.js                 # SIWE config for AppKit
│   ├── components/                 # Reusable UI components
│   │   ├── NavBar.vue              # Navigation bar with wallet connect button
│   │   ├── TheWelcome.vue          # Landing page hero section
│   │   └── icons/                  # Icon helper stubs
│   ├── config/
│   │   └── common.js               # Hardcoded config: BACKEND_URL, ACCESS_TOKEN_KEY
│   ├── router/
│   │   └── index.js                # Vue Router (hash history, lazy-loaded routes)
│   ├── stores/
│   │   └── auth.js                 # Pinia store (token state: setToken, clearToken)
│   └── views/
│       ├── HomeView.vue            # Landing page
│       └── AccountView.vue         # Account page
├── public/                         # Static assets
├── vite.config.js                  # Vite build config (+ Vue plugin, `@` alias → ./src)
├── vitest.config.js                # Testing config (jsdom environment)
├── eslint.config.js                # ESLint flat config
├── .oxlintrc.json                  # Oxlint rules
├── .prettierrc.json                # Prettier: semi=false, singleQuote=true, printWidth=100
├── jsconfig.json                   # Path alias `@` → ./src
├── package.json
└── index.html
```

## Key Patterns

### Routing
- Uses **vue-router v5** with `createWebHistory` mode.
- Routes are defined in `src/router/index.js`.
- Non-primary routes should be **lazy-loaded** via dynamic `import()`.

### State Management
- Uses **Pinia** with composition API style stores.
- Auth token lives in `src/stores/auth.js` (`useAuthStore`).
- Keep stores lean; avoid duplicating backend state locally.

### Authentication Flow (SIWE + JWT)
1. User clicks connect wallet → AppKit prompts wallet signature.
2. `getNonce()` fetches a fresh nonce from backend (`GET /api/v1/auth/nonce`).
3. User signs SIWE message → `verifyMessage()` POSTs signature to backend (`POST /api/v1/auth/verify`).
4. On success, backend returns JWT → stored via `useAuthStore.setToken()`.
5. Subsequent requests include `Authorization: Bearer <token>` header.
6. Token is refreshed on reconnect via `GET /api/v1/auth/refresh`.

### API / Service Layer
- Currently uses raw `fetch()` calls directly in `src/auth/authenticate.js`.
- No centralized HTTP client abstraction yet — this is a known technical debt area.
- All API endpoints are under `http://localhost:3000` (hardcoded in `src/config/common.js`).
- Known bug: template literal error strings in `authenticate.js` use single quotes instead of backticks (lines 27, 42, 46, 69).

### Component Conventions
- All components use `<script setup>` Composition API syntax.
- Use Bootstrap utility classes for styling where possible.
- Keep components focused and small; split into sub-components as needed.

## Commands

```bash
npm run dev        # Start Vite dev server
npm run build      # Production build to dist/
npm run preview    # Preview production build locally
npm run test:unit  # Run Vitest unit tests
npm run lint       # Run oxlint + eslint sequentially
npm run lint:ox    # Run oxlint only
npm run lint:eslint # Run ESLint with auto-fix
npm run format     # Run Prettier on src/
```

## Testing

- Config: `vitest.config.js` (jsdom environment, merges with Vite config).
- Test files go in `src/**/__tests__/*.test.{js,jsx}`.
- Use `@vue/test-utils` for component testing.
- No tests exist yet — new features and bug fixes should include tests.

## Coding Guidelines

- Follow the existing code style: no semicolons, single quotes, max line width 100 (Prettier).
- Import order: Vue/Pinia/router first, then third-party, then local modules (`@/` alias).
- Use kebab-case for HTML elements/attributes, camelCase for JS/Vue APIs.
- Write defensively: always check HTTP response status before parsing JSON body.
- Never commit secrets (project IDs, private keys, tokens).
- When adding new API endpoints, document them alongside the existing auth calls in `src/auth/authenticate.js` or extract a proper service layer.

## Known Issues & Technical Debt

1. **Hardcoded values**: `BACKEND_URL` and Reown `PROJECT_ID` should be moved to `.env` variables (`VITE_BACKEND_URL`, `VITE_PROJECT_ID`).
2. **Template literal bugs** in `src/auth/authenticate.js`: single-quoted strings on lines 27, 42, 46, 69 use `${}` interpolation but won't interpolate.
3. **Mixed typing**: `AccountView.vue` declares `<script setup lang="ts">` but contains no TypeScript — remove `lang="ts"` or convert properly.
4. **No centralized API layer**: consider extracting a reusable HTTP client/utility.
5. **No tests**: testing infrastructure is ready but unpopulated.
6. **Stale scaffolding**: `index.html` title says "Vite App", some icon components are stubs.
