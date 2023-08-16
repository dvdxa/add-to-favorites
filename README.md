# add-to-favorites
API allows users who monitor terminals, to add/remove terminals to/from favorites and simplify their work. They don't need to look for necessary terminals by scrolling up and down. Favorited terminals will be displayed at the top. When user removes terminal from favorites, it will be returned back to default terminals list and will displayed in default order.

## Installation
To install this App, follow these steps:
1. Clone the repository
2. Install required dependencies
3. Set up the database(postgres) and configure it

## Usage
To use this App:
1. Explore endpoints in `./internal/handlers/ports.go`
2. Run the app: `go run ./cmd/main.go`
3. Open desired testing tool such as Postman or Insomnia
4. Sign up for an account or log in if you already have one
5. Add/remove terminals to/from favorites

