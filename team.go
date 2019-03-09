package gotournament

// TeamInterface defines the methods for teams. Teams are used to create tournaments and generate games.
// Teams can have games and can contain a slice of players
type TeamInterface interface {
	GetID() int
	GetPlayers() []PlayerInterface
	GetGames() []GameInterface
	AppendGame(game GameInterface)
}

// Team is a default struct used as an example of how structs can be implemented for gotournament
type Team struct {
	ID      int
	Players []PlayerInterface
	Games   []GameInterface
}

// GetID returns the id of the score
func (t *Team) GetID() int {
	return t.ID
}

// GetPlayers returns a slice of players
func (t *Team) GetPlayers() []PlayerInterface {
	return t.Players
}

// GetGames returns a slice of games
func (t *Team) GetGames() []GameInterface {
	return t.Games
}

// AppendGame takes a game as an argument and appends it to the Games slice
func (t *Team) AppendGame(game GameInterface) {
	t.Games = append(t.Games, game)
}
