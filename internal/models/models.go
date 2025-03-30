package models


type Vote struct {
	ID       string            
	ChannelID string           
	Options  []string          
	Results  map[string]int    
	Active   bool              
}


type VoteOption struct {
	Name  string
	Count int
}