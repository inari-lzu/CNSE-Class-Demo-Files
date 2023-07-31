package voterList

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type voterPoll struct {
	PollID   uint      `json:"id"`
	VoteDate time.Time `json:"date"`
}

type Voter struct {
	VoterID     uint        `json:"id"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	VoteHistory []voterPoll `json:"history"`
}

type VoterList struct {
	Voters map[uint]Voter `json:"voters"`
}

func NewVoterList() (*VoterList, error) {
	voteList := &VoterList{
		Voters: make(map[uint]Voter),
	}
	return voteList, nil
}

func (v *VoterList) AddVoter(voter Voter) error {
	_, ok := v.Voters[voter.VoterID]
	if ok {
		return errors.New("voter already exists")
	}
	v.Voters[voter.VoterID] = voter
	return nil
}

func (v *VoterList) UpdateVoter(voter Voter) (Voter, error) {
	vt, ok := v.Voters[voter.VoterID]
	if !ok {
		return Voter{}, errors.New("voter does not exist")
	}
	voter.VoteHistory = vt.VoteHistory // only modify resource data
	v.Voters[voter.VoterID] = voter
	return voter, nil
}

func (v *VoterList) DeleteVoter(id uint) error {
	_, ok := v.Voters[id]
	if !ok {
		return errors.New("voter does not exist")
	}
	delete(v.Voters, id)
	return nil
}

func (v *VoterList) Clear() error {
	v.Voters = make(map[uint]Voter)
	return nil
}

func (v *VoterList) GetVoter(id uint) (Voter, error) {
	voter, ok := v.Voters[id]
	if !ok {
		return Voter{}, errors.New("voter does not exist")
	}
	return voter, nil
}

func (v *VoterList) GetVoterPoll(voterId uint, pollid uint) (voterPoll, error) {
	voter, ok := v.Voters[voterId]
	if !ok {
		return voterPoll{}, errors.New("voter does not exist")
	}
	for _, vp := range voter.VoteHistory {
		if vp.PollID == pollid {
			return vp, nil
		}
	}
	return voterPoll{}, errors.New("voter poll does not exist")
}

func (v *VoterList) AddVoterPoll(voterId uint, newvp voterPoll) error {
	voter, ok := v.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}
	for _, vp := range voter.VoteHistory {
		if newvp.PollID == vp.PollID {
			return errors.New("voter poll already exists")
		}
	}
	voter.VoteHistory = append(voter.VoteHistory, newvp)
	v.Voters[voterId] = voter
	return nil
}

func (v *VoterList) UpdateVoterPoll(voterId uint, newvp voterPoll) error {
	voter, ok := v.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}
	for i, vp := range voter.VoteHistory {
		if newvp.PollID == vp.PollID {
			voter.VoteHistory[i] = newvp
			v.Voters[voterId] = voter
			return nil
		}
	}
	return errors.New("voter poll does not exist")
}

func (v *VoterList) DeleteVoterPoll(voterId uint, pollid uint) error {
	voter, ok := v.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}
	for i, vp := range voter.VoteHistory {
		if vp.PollID == pollid {
			voter.VoteHistory = append(voter.VoteHistory[:i], voter.VoteHistory[i+1:]...)
			v.Voters[voterId] = voter
			return nil
		}
	}
	return errors.New("voter poll does not exist")
}

func (v *VoterList) GetAllVoters() ([]Voter, error) {
	var vl []Voter
	for _, voter := range v.Voters {
		vl = append(vl, voter)
	}
	return vl, nil
}

func (v *VoterList) PrintVoter(voter Voter) {
	jsonBytes, _ := json.MarshalIndent(voter, "", "  ")
	fmt.Println(string(jsonBytes))
}

func (v *VoterList) PrintAllVoters(vl []Voter) {
	for _, voter := range vl {
		v.PrintVoter(voter)
	}
}

func (v *VoterList) JsonToVoter(jsonString string) (Voter, error) {
	var voter Voter
	err := json.Unmarshal([]byte(jsonString), &voter)
	if err != nil {
		return Voter{}, err
	}

	return voter, nil
}
