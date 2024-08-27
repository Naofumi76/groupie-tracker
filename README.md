# Groupie Tracker

Groupie Tracker is a web application designed to display information about musical artists. Utilizing both back-end and front-end technologies, the application fetches data from an external API to present details such as artist images, group members, and other relevant information. Users can view these details either on the main page or by navigating to individual artist pages.

## Features:

Changes from Groupie-Tracker to Groupie-Tracker-Filters:

- Added major filters that affect the display of artists on the main page.

- When no filters are selected or on initial page load, all artists are displayed.

- Depending on the selected filters, the displayed artists will vary:

**I) Members:**
- Artists are displayed based on the number of members selected. Multiple checkboxes can be selected to display groups that match any selected member count.

**II) Creation Date:**
- Artists are displayed if their creation date matches the selected value. Selecting the "After" checkbox includes artists with creation dates on or after the selected value.

**III) Locations:**
- Artists are displayed if they have performed concerts in locations that match any selected location. Locations are capitalized and formatted for readability.

**IV) First Album Date:**
- Artists are displayed if their first album date matches the selected value. Selecting the "After" checkbox includes artists with first album dates on or after the selected value.

[IMPORTANT] Please use the "Filter" button to apply your selected filters.


## Technologies Used:

    Front-end: HTML, CSS
    Back-end: Golang
    Templates: Go templates for rendering HTML
    External API: The project uses a pre-existing API for artist data.

# Project Setup

## Clone the repository:

```bash
git clone https://github.com/cnuttens/groupie-tracker-filters.git
```

## Navigate to the project directory:


```bash
cd groupie-tracker
```

## Run the project:


```bash
go run .
```
## Usage:

    Access the application in your web browser at http://localhost:8080.
    Navigate through the main page to see a list of artists.
    Click on an artist to view detailed information on a new page.


## API Information:

    The application fetches data from https://groupietrackers.herokuapp.com/api/artists.
    Relation data is fetched from https://groupietrackers.herokuapp.com/api/relation.