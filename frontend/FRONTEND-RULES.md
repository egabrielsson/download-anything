# React + Tailwind Best Practices (Basics for Starters)

Focus on clean components, custom hooks, and Tailwind usage for a maintainable frontend.

## Project Structure (Feature-based for scalability)

src/
├── assets/                  # Images, fonts
├── components/              # Reusable UI (Button, Card, etc.)
├── features/                # Domain/feature folders (e.g., scraper/)
│   └── scraper/
│       ├── components/      # Feature-specific components
│       ├── hooks/           # Feature-specific hooks
│       └── pages/           # Page-level for this feature
├── hooks/                   # Global custom hooks (useFetch, useAuth)
├── pages/                   # Route pages (Home, Dashboard)
├── layouts/                 # Layout wrappers (MainLayout)
├── services/                # API calls (axios/fetch wrappers)
├── utils/                   # Helpers (formatDate, cn for classes)
├── App.tsx
└── main.tsx

- Start with `components/` + `hooks/` + `pages/` for small apps.
- Move to `features/` when you have multiple domains (e.g., auth, scraper, results).

## Components
- Make small, reusable components.
- Use composition: prefer children props over props drilling.
- One component per file; export default.
- Use TypeScript for props (interface Props { ... }).

## Custom Hooks Abstraction
- Extract reusable logic into hooks (e.g., `useScraperData()` for fetching/parsing).
- Follow rules of hooks: call at top level, only in components/hooks.
- Examples: `useFetch`, `useForm`, `useLocalStorage`.

## Routing
- Use React Router (v6+): `<BrowserRouter>`, `<Routes>`, `<Route>`.
- Colocate routes with pages/features (e.g., `pages/Dashboard.tsx`).
- Use loaders/actions if using data APIs (or TanStack Query).

## Tailwind CSS Tips
- Use full class names (no concatenation): `bg-blue-500`, not `bg-${color}-500`.
- Conditional classes: use `cn` helper (from `clsx` + `tailwind-merge`):
  ```tsx
  import { cn } from '@/utils/cn';

  <div className={cn("p-4", isActive && "bg-blue-100")} />

  Mobile-first: text-sm md:text-base lg:text-lg.
Dark mode: dark:bg-gray-800.
Group utilities: group-hover:scale-105.
Install prettier-plugin-tailwindcss to auto-sort classes.
Avoid overusing arbitrary values [...]; prefer theme extensions if needed.

General React Tips

Use functional components + hooks (no class components).
Manage state locally first; lift up when shared.
Use TanStack Query or SWR for data fetching (better than useEffect + fetch).
Memoize expensive computations (useMemo, useCallback).
Keep components pure and predictable.

Start small, extract when things repeat. These scale well for a web scraper UI.


<grok-card data-id="f6d98e" data-type="image_card" data-plain-type="render_searched_image"  data-arg-size="LARGE" ></grok-card>

<grok-card data-id="39f2de" data-type="image_card" data-plain-type="render_searched_image"  data-arg-size="LARGE" ></grok-card>

(Example folder structure visuals—yours can be simpler at first.)

These should give you a solid, no-frills starting point. Adjust as your vibecode project grows! If you need examples or expansions later, let me know.