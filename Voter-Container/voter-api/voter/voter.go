package voter

import (
	"errors"
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

func (v Voter) GetID() uint {
	return v.VoterID
}

func NewVoter() Voter {
	v := Voter{}
	v.VoteHistory = make([]voterPoll, 0)
	return v
}

func (v Voter) UpdateVoterInfo(vt Voter) (Voter, error) {
	vt.VoterID = v.VoterID
	vt.VoteHistory = v.VoteHistory // we don't expect to voteHistory
	return vt, nil
}

func (v Voter) UpdateVoteHistory(vt Voter) (Voter, error) {
	v.VoteHistory = vt.VoteHistory
	return v, nil
}

func (v Voter) GetHistory() []voterPoll {
	return v.VoteHistory
}

func (v Voter) DeleteHistory() Voter {
	v.VoteHistory = make([]voterPoll, 0)
	return v
}

func (v Voter) AddPoll(newvp voterPoll) (Voter, error) {
	for _, vp := range v.VoteHistory {
		if newvp.PollID == vp.PollID {
			return Voter{}, errors.New("v poll already exists")
		}
	}
	v.VoteHistory = append(v.VoteHistory, newvp)
	return v, nil
}

func (v Voter) GetPoll(pollid uint) (voterPoll, error) {
	for _, vp := range v.VoteHistory {
		if vp.PollID == pollid {
			return vp, nil
		}
	}
	return voterPoll{}, errors.New("v poll does not exist")
}

func (v Voter) UpdatePoll(newvp voterPoll) (Voter, error) {
	pollid := newvp.PollID
	for i, vp := range v.VoteHistory {
		if vp.PollID == pollid {
			v.VoteHistory[i] = newvp
			return v, nil
		}
	}
	return Voter{}, errors.New("v poll does not exist")
}

func (v Voter) DeletePoll(pollid uint) (Voter, error) {
	for i, vp := range v.VoteHistory {
		if vp.PollID == pollid {
			v.VoteHistory = append(v.VoteHistory[:i], v.VoteHistory[i+1:]...)
			return v, nil
		}
	}
	return Voter{}, errors.New("v poll does not exist")
}
