## Aman Swami (Backend)




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