# Recipe App TODO List

This document outlines the tasks needed to implement the feature roadmap for the Recipe App.

---

## Phase 1: Core MVP – Passive & Functional Essentials

### Backend Services
- [ ] **User Management**
  - [ ] Create User CRUD endpoints (create, read, update, delete)
  - [ ] Implement authentication (login, logout, password reset)
  - [ ] Set up user profile management (dietary preferences, favorites)
- [ ] **Recipe Repository & Modifications**
  - [ ] Integrate pre-generated 10k recipes into the database
  - [ ] Create endpoints for retrieving recipes
  - [ ] Build endpoints for fetching basic modifications for each recipe
  - [ ] Implement the RAG (Retrieval-Augmented Generation) system on AWS for quick content retrieval
- [ ] **Meal Planning**
  - [ ] Develop endpoints for saving weekly meal plans
  - [ ] Create endpoints to generate simple shopping lists from selected recipes
- [ ] **Analytics & Data Collection**
  - [ ] Set up basic logging for user interactions (recipe views, meal planning usage)
  - [ ] Integrate analytics to track app usage and engagement

### Frontend / App UI
- [ ] **User Interface**
  - [ ] Design and implement a clean UI for browsing recipes
  - [ ] Build a basic search and filter interface (by ingredients, cuisine, dietary restrictions)
- [ ] **User Authentication & Accounts**
  - [ ] Create signup/login screens
  - [ ] Implement profile management UI
- [ ] **Meal Planning Interface**
  - [ ] Develop a weekly meal planner view
  - [ ] Implement functionality to save and display meal plans
  - [ ] Create a view for generating a basic shopping list

### Infrastructure & Early Access
- [ ] **Hosting & Deployment**
  - [ ] Set up AWS instance to run the backend services (RAG system and app logic)
  - [ ] Configure domain and DNS settings
- [ ] **User Testing Setup**
  - [ ] Create a system for early access (first 50–100 testers with unlimited, ad‑free access)
  - [ ] Set up feedback channels (in-app feedback form, email, etc.)

---

## Phase 2: Enhanced Engagement & Initial Active Features

### Backend Services
- [ ] **Advanced Meal Planning**
  - [ ] Expand meal planning endpoints to include customization (nutritional goals, detailed menus)
  - [ ] Implement endpoints for saving custom meal plans and recipe history
- [ ] **Manual Inventory Tracking**
  - [ ] Create endpoints to allow users to manually add items to their pantry/fridge
  - [ ] Implement logic for two modes:
    - [ ] **Recipe Generation Mode:** Suggest recipes based on current inventory
    - [ ] **Shopping List Mode:** Generate a list of missing items for a recipe
  - [ ] Build functionality to update inventory after a recipe is confirmed as made
  - [ ] **Barcode/Nutritional Label Scanning Integration:**
    - [ ] Develop API to process barcode or label data and add items to inventory
- [ ] **Grocery Integration (Initial)**
  - [ ] Develop endpoints to integrate with one or two grocery/delivery APIs (e.g., Instacart)
  - [ ] Create service to convert meal plans into grocery orders or optimized shopping lists
- [ ] **Social & Community Features (Basic)**
  - [ ] Create endpoints for recipe ratings and comments
  - [ ] Set up endpoints for sharing recipes on social media
- [ ] **Monetization Foundations**
  - [ ] Implement endpoints for microtransactions (purchase extra recipes/modifications)
  - [ ] Set up subscription tier endpoints (Basic, Standard)
  - [ ] Integrate minimal in-app ad serving endpoints for free users

### Frontend / App UI
- [ ] **Advanced Meal Planning UI**
  - [ ] Update meal planner to include customization options (nutritional goals, detailed menus)
  - [ ] Build interface to save and view custom meal plans and history
- [ ] **Inventory Tracking UI**
  - [ ] Develop a UI for manual inventory input (pantry/fridge items)
  - [ ] Integrate barcode scanning functionality in the app for quick item addition
  - [ ] Create views for recipe generation mode and shopping list mode based on inventory
- [ ] **Grocery Integration UI**
  - [ ] Implement a screen to view and manage grocery orders/shopping lists from meal plans
- [ ] **Social & Community UI**
  - [ ] Add UI elements for rating and commenting on recipes
  - [ ] Provide social sharing buttons/options on recipe pages
- [ ] **Monetization UI**
  - [ ] Design subscription signup pages for each tier (Basic, Standard)
  - [ ] Integrate microtransaction options within the recipe/modification screens
  - [ ] Implement basic ad placements for free-tier users

### Testing & Deployment
- [ ] **Integration Testing**
  - [ ] Test all new endpoints and UI flows with inventory, grocery integration, and meal planning enhancements
- [ ] **User Feedback**
  - [ ] Deploy beta for a broader user group and gather usage data and feedback for Phase 2 features
- [ ] **Documentation**
  - [ ] Update API documentation and user guides as new features are implemented

---

*Note:* This TODO list is meant to be a living document. Adjust priorities, add more granular tasks, or split tasks further as the project evolves.
