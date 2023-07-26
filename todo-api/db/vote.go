package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type VoterPoll struct {
	PollID   uint      `json:"id"`
	VoteDate time.Time `json:"date"`
}

type Voter struct {
	VoterID     uint        `json:"id"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	VoteHistory []VoterPoll `json:"history"`
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

func (vl *VoterList) AddVoter(voter Voter) error {
	_, ok := vl.Voters[voter.VoterID]
	if ok {
		return errors.New("voter already exists")
	}
	vl.Voters[voter.VoterID] = voter
	return nil
}

func (vl *VoterList) UpdateVoter(voter Voter) (Voter, error) {
	v, ok := vl.Voters[voter.VoterID]
	if !ok {
		return Voter{}, errors.New("voter does not exist")
	}
	voter.VoteHistory = v.VoteHistory // only modify resource data
	vl.Voters[voter.VoterID] = voter
	return voter, nil
}

func (vl *VoterList) DeleteVoter(id uint) error {
	_, ok := vl.Voters[id]
	if !ok {
		return errors.New("voter does not exist")
	}
	delete(vl.Voters, id)
	return nil
}

func (vl *VoterList) Clear() error {
	vl.Voters = make(map[uint]Voter)
	return nil
}

func (vl *VoterList) GetVoter(id uint) (Voter, error) {
	voter, ok := vl.Voters[id]
	if !ok {
		return Voter{}, errors.New("voter does not exist")
	}
	return voter, nil
}

func (vl *VoterList) GetVoterPoll(voterId uint, pollid uint) (VoterPoll, error) {
	voter, ok := vl.Voters[voterId]
	if !ok {
		return VoterPoll{}, errors.New("voter does not exist")
	}
	for _, voterPoll := range voter.VoteHistory {
		if voterPoll.PollID == pollid {
			return voterPoll, nil
		}
	}
	return VoterPoll{}, errors.New("voter poll does not exist")
}

func (vl *VoterList) AddVoterPoll(voterId uint, voterPoll VoterPoll) error {
	voter, ok := vl.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}
	for _, vp := range voter.VoteHistory {
		if voterPoll.PollID == vp.PollID {
			return errors.New("voter poll already exists")
		}
	}
	voter.VoteHistory = append(voter.VoteHistory, voterPoll)
	vl.Voters[voterId] = voter
	return nil
}

func (vl *VoterList) UpdateVoterPoll(voterId uint, voterPoll VoterPoll) error {
	voter, ok := vl.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}
	for i, vp := range voter.VoteHistory {
		if voterPoll.PollID == vp.PollID {
			voter.VoteHistory[i] = voterPoll
			vl.Voters[voterId] = voter
			return nil
		}
	}
	return errors.New("voter poll does not exist")
}

func (vl *VoterList) DeleteVoterPoll(voterId uint, pollid uint) error {
	voter, ok := vl.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}
	for i, voterPoll := range voter.VoteHistory {
		if voterPoll.PollID == pollid {
			voter.VoteHistory = append(voter.VoteHistory[:i], voter.VoteHistory[i+1:]...)
			vl.Voters[voterId] = voter
			return nil
		}
	}
	return errors.New("voter poll does not exist")
}

func (vl *VoterList) GetAllVoters() ([]Voter, error) {
	var voterList []Voter
	for _, voter := range vl.Voters {
		voterList = append(voterList, voter)
	}
	return voterList, nil
}

func (vl *VoterList) PrintVoter(voter Voter) {
	jsonBytes, _ := json.MarshalIndent(voter, "", "  ")
	fmt.Println(string(jsonBytes))
}

func (vl *VoterList) PrintAllVoters(voterList []Voter) {
	for _, voter := range voterList {
		vl.PrintVoter(voter)
	}
}

func (vl *VoterList) JsonToVoter(jsonString string) (Voter, error) {
	var voter Voter
	err := json.Unmarshal([]byte(jsonString), &voter)
	if err != nil {
		return Voter{}, err
	}

	return voter, nil
}
