package gotournament

// TournamentInterface defines the methods needed to handle tournaments.
type TournamentInterface interface {
	GetType() int
	GetTeams() []TeamInterface
	GetGroups() []TournamentGroupInterface
	GetGames() []GameInterface
}

// Tournament is a default struct used as an example of how structs can be implemented for gotournament
type Tournament struct {
	Type   TournamentType // Is it elimination or group or ladder or poker? What is a type?
	Teams  []TeamInterface
	Groups []TournamentGroupInterface
	Games  []GameInterface
}

// GetType returns the type of tournament
func (t Tournament) GetType() int {
	return int(t.Type)
}

// GetTeams returns the team slice
func (t Tournament) GetTeams() []TeamInterface {
	return t.Teams
}

// GetGroups returns the group slice
func (t Tournament) GetGroups() []TournamentGroupInterface {
	return t.Groups
}

// GetGames returns the game slice
func (t Tournament) GetGames() []GameInterface {
	return t.Games
}

// CreateTournament creates a tournament with the simplest input. It is recommended to create a slice with
// specific use via CreateTournamentFromTeams as this method will generate it's own Teams as a sort of placeholder.
func CreateTournament(teamCount int, meetCount int, tournamentType int) TournamentInterface {
	var teams []TeamInterface

	for i := 0; i < teamCount; i++ {
		teams = append(teams, &Team{ID: i})
	}

	return CreateTournamentFromTeams(teams, meetCount, tournamentType)
}

// CreateTournamentFromTeams takes a slice of teams and generates a tournament of the specified type
func CreateTournamentFromTeams(teams []TeamInterface, meetCount int, tournamentType int) TournamentInterface {
	if TournamentType(tournamentType) == TournamentTypeGroup {
		// It is recommended to call CreateGroupTournamentFromTeams directly as we try to automatically determine a group size and count here
		groupCount := len(teams) / 4 // We assume 4 teams per group by default
		if groupCount < 1 {
			groupCount = 1
		}
		return CreateGroupTournamentFromTeams(teams, groupCount, meetCount)
	}
	return nil
}

// CreateGroupTournamentFromTeams takes a slice of teams and generates a group tournament
func CreateGroupTournamentFromTeams(teams []TeamInterface, groupCount int, meetCount int) TournamentInterface {
	var groups []TournamentGroupInterface

	for i := 0; i < groupCount; i++ {
		groups = append(groups, &TournamentGroup{ID: i})
	}

	for i, team := range teams {
		groupIndex := i % (len(groups))
		groups[groupIndex].AppendTeam(team)
	}

	return CreateGroupTournamentFromGroups(groups, meetCount)
}

// CreateGroupTournamentFromGroups takes a slice of groups that contain teams and returns a group tournament
func CreateGroupTournamentFromGroups(groups []TournamentGroupInterface, meetCount int) TournamentInterface {
	// Works best for an even amount of teams in every group
	var games []GameInterface
	var teams []TeamInterface
	gameIndex := 0
	for gi, group := range groups {
		groupteams := group.GetTeams()
		if len(groupteams) < 4 {
			// TODO remove panic
			panic("Group must contain at least 4 teams, currently" + string(len(groupteams)))
		}
		teams = append(teams, group.GetTeams()...)
		// Loop through meet count
		for mi := 0; mi < meetCount; mi++ {
			var homeTeams []TeamInterface
			var awayTeams []TeamInterface
			// Everyone meets everyone once
			// We begin by taking our array of teams like 0,1,2,3, and splitting it
			if mi%2 == 0 {
				homeTeams = group.GetTeams()[0:(len(group.GetTeams()) / 2)]
				awayTeams = group.GetTeams()[(len(group.GetTeams()) / 2):]
			} else {
				awayTeams = group.GetTeams()[0:(len(group.GetTeams()) / 2)]
				homeTeams = group.GetTeams()[(len(group.GetTeams()) / 2):]
			}
			for i := 0; i < len(group.GetTeams())-1; i++ {
				// Now we have home teams of 0,1 and away teams of 2,3
				// This means 0 will meet 2 and 1 will meet 3
				for hi, hteam := range homeTeams {
					game := Game{Teams: []TeamInterface{hteam, awayTeams[hi]}}
					groups[gi].AppendGame(&game)
					hteam.AppendGame(&game)
					games = append(games, &game)
					awayTeams[hi].AppendGame(&game)
					gameIndex++
				}
				var x, y, z TeamInterface
				// We keep the first home team in the same position and rotate all others
				x, homeTeams = homeTeams[0], homeTeams[1:]
				// Take the first away team
				z, awayTeams = awayTeams[0], awayTeams[1:]
				// and append to end of home teams
				homeTeams = append(homeTeams, z)
				// Take the second home team
				y, homeTeams = homeTeams[0], homeTeams[1:]
				// and append it to the end of away teams
				awayTeams = append(awayTeams, y)
				// Put the first home team back in first position of home array
				homeTeams = append([]TeamInterface{x}, homeTeams...)
			}
		}
	}
	return Tournament{Groups: groups, Games: games, Teams: teams, Type: TournamentTypeGroup}
}

// NumberOfGamesForGroupTournament Calculates the number of games in a group tournament based on number of teams, groups and unique encounters.
func NumberOfGamesForGroupTournament(teamCount int, groupCount int, meetCount int) int {
	return ((((teamCount / groupCount) - 1) * ((teamCount / groupCount) / 2)) * groupCount) * meetCount
}
