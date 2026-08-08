## Aman Swami (Backend)

1. Create an autonomous AI persona that can be initialized once and then continue operating without further human input.

2. Implement `POST /api/agent/init` so the frontend can send a persona payload and receive a stable `agentId` in response.

3. Implement `GET /api/agent/feed?agentId=...` so the frontend can retrieve a live list of published posts for the agent.

4. Build topic discovery logic that pulls from live sources and presents only the most compelling candidates to the agent.

5. Give the agent a strong editorial voice, clear domain focus, and a persona that remains consistent across published posts.

6. Make the agent remember what it has already published, so it avoids repeating the same ideas or recycled topic summaries.

7. Add a scheduler that publishes new posts gradually over time, instead of publishing all content at once.

8. Attach rationale, relevance, and source details to every published post so the output is explainable and traceable.

9. Store all agent state in PostgreSQL so published posts, persona metadata, and seen topics survive restarts.

10. Create frontend-friendly dashboard endpoints for agent details, activity, topics, memory, analytics, and logs.

11. Ensure response payloads follow a consistent JSON contract and always use ISO 8601 timestamps for dates.

12. Validate the `agentId` parameter on every request and return a clear error if it is missing or invalid.

13. Return `success: true` or `success: false` where appropriate to keep the API responses predictable.

14. Keep the agent’s `name` and `domain` stored in the backend so the frontend can show persona metadata in the UI.

15. Use consistent field names across all endpoints so the frontend does not need special-case response parsing.

16. Guarantee the feed returns an empty array rather than `null` when no posts exist yet.

17. Persist `seenTopicIDs` so the agent does not re-evaluate the same topic more than once.

18. Use a simple backend entity model with Agent, Persona, Post, Topic, and memory state.

19. Provide a startup path that works even if PostgreSQL is temporarily unavailable, with an in-memory fallback during development.

20. Always include a source URL in each post’s `sources` array, even if it is a placeholder from the discovery pipeline.

21. Return agent activity metadata such as `ACTIVE`, `Searching`, `Online`, and the current task progress.

22. Keep the backend API contract aligned with React dashboard routes so the frontend can integrate smoothly.

23. Ensure the scheduler starts immediately after agent initialization and runs on a repeat interval.

24. Keep the init endpoint fast so the frontend receives the ID immediately while the agent continues its work asynchronously.

25. Allow the frontend to poll live feed updates every few seconds without overloading the backend.

26. Keep the database schema minimal and reliable so the system is easy to maintain during the hackathon.

27. Provide a fallback in-memory state path if the database schema cannot load or the connection fails.

28. Add regression tests that verify API responses, endpoint shapes, and required contract fields.

29. Return `createdAt` timestamps in RFC3339 / ISO 8601 format for all posts and metadata.

30. Keep error responses clear, with helpful messages for missing parameters and invalid IDs.

31. Design initialization to be idempotent enough that repeated client retries do not create duplicate agents.

32. Make the feed endpoint the canonical source of truth for all published post data.

33. Keep the backend routes simple and easy to route from the React frontend.

34. Give every post a unique ID and preserve that ID in the feed output.

35. Provide a topic queue endpoint that can show accepted/rejected subjects, even before publishing begins.

36. Keep the agent publishing flow transparent by returning rationale and sources in every feed item.

37. Document the backend setup, how to run the server locally, and how to start the frontend in `backend/README.md`.

38. Keep the agent’s editorial standard strict: no fluff, no recycled takes, and no weak topic choices.

39. Make the backend able to store agent memory and retrieve it through a dedicated memory endpoint.

40. Keep the API endpoints robust enough for both development and the hackathon evaluator.



## Mohit Sharma (Frontend)
1. Initialize the frontend application.

    Tasks:
    1. Verify the existing React + Vite setup.
    2. Install and configure Tailwind CSS.
    3. Install the required frontend dependencies:
    - react-router-dom
    - axios
    - framer-motion
    - lucide-react
    4. Remove the default Vite starter code.
    5. Create a clean starter application with a minimal App component.
    6. Ensure the project builds and runs successfully.
    7. Keep the project clean and production-ready.
    8. Work only inside the frontend directory.
    9. After completion, summarize the changes made.

2. The starter UI was only for setup verification.

    Remove all starter/demo UI and leave a clean React app.

    Keep the project configuration intact (Tailwind, Router, Axios, Framer Motion, Lucide).

    Render only a centered "Frontend Initialized" placeholder in App.jsx.

    Don't create pages, components, or features yet.

    After finishing:
    - verify the app runs successfully
    - add this exact prompt under my section in PROMPTS.md using the next prompt number

3. PROMPTS.md Logging Policy

    This policy applies to every future task unless I explicitly disable it.

    - Before implementing any task that results in meaningful code changes, automatically log my EXACT prompt in PROMPTS.md.
    - Append it only under "## Mohit Sharma (Frontend)".
    - Automatically determine the next sequential prompt number.
    - Never modify, renumber, or delete previous prompt entries.
    - Never modify my teammate's section.
    - Do not log prompts that are only questions, discussions, or planning conversations.

4. Configure React Router for the application.

    Create only two pages:
    - Landing
    - Dashboard
    Dashboard should accept :agentId.
    Keep routing clean and simple.

    Theme
    Background- #222831
    Surface- #393E46
    Primary-#948979
    Text- #DFD0B8

    Design Style
    - Modern SaaS
    - Minimal
    - Dark Theme
    - Rounded Corners
    - Smooth Animations
    - Desktop First

5. Build the Landing page.

    The page should allow the user to:
    - Enter Persona Name
    - Select AI Domain
    - Click Initialize Agent
    Don't connect the backend yet.

6. Create DashboardLayout.
    It should contain:
    - Sidebar
    - Navbar
    - Main Content Area
    Use React Router Outlet.
    Only build the layout.

7. Create reusable UI components:
   - Button
   - Card
   - Input
   - Badge
   - Loading Spinner
   Keep styling consistent with the project theme.

8. Build the dashboard navbar.
    Include:
    - Agent Name
    - Agent Status
    - Theme-ready layout
9. Build the Dashboard overview.
    Include cards for:
    - Current Activity
    - Statistics
    - Recent Feed
    - Topic Queue
    - Memory Summary
    Use mock data.
10. Build the Feed section.
    Display posts with:
    - Content
    - Time
    - Rationale
    - Sources
    Use mock data.
11. Build the Topics section.
    Display:
    Accepted Topics
    Rejected Topics
    Reason for rejection
    Use mock data.
12. Build the Memory section.
    Display:
    - Interests
    - Recent Topics
    Use mock data.
13. Build the Analytics section.
    Display statistics returned by the backend.
    Use mock data.
14. Replace mock data with API calls using Axios.
    Use the API contract.
    Keep the UI unchanged.
