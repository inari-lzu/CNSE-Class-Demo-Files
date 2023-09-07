package voter

import (
	"errors"
	"strconv"
	"time"
)

type voterPoll struct {
	PollID   uint      `json:"id"`
	VoteDate time.Time `json:"date"`
}

type voterPollJson struct {
	VoterPoll string `json:"voterPoll"`
	Date      string `json:"date"`
}

type Voter struct {
	VoterID     uint        `json:"id"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	VoteHistory []voterPoll `json:"history,omitempty"`
}

type VoterJson struct {
	Voter     string   `json:"voter"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	History   []string `json:"history,omitempty"`
}

func (vr Voter) GetID() uint {
	return vr.VoterID
}

func NewVoter() Voter {
	vr := Voter{}
	vr.VoteHistory = make([]voterPoll, 0)
	return vr
}

func (vr Voter) UpdateVoterInfo(vt Voter) (Voter, error) {
	vt.VoterID = vr.VoterID
	vt.VoteHistory = vr.VoteHistory // we don't expect to change voteHistory in updating
	return vt, nil
}

func (vr Voter) UpdateVoteHistory(vt Voter) (Voter, error) {
	vr.VoteHistory = vt.VoteHistory
	return vr, nil
}

func (vr Voter) GetHistory() []voterPoll {
	return vr.VoteHistory
}

func (vr Voter) DeleteHistory() Voter {
	vr.VoteHistory = make([]voterPoll, 0)
	return vr
}

func (vr Voter) AddPoll(newvp voterPoll) (Voter, error) {
	for _, vp := range vr.VoteHistory {
		if newvp.PollID == vp.PollID {
			return Voter{}, errors.New("vr poll already exists")
		}
	}
	vr.VoteHistory = append(vr.VoteHistory, newvp)
	return vr, nil
}

func (vr Voter) GetPoll(pollid uint) (voterPoll, error) {
	for _, vp := range vr.VoteHistory {
		if vp.PollID == pollid {
			return vp, nil
		}
	}
	return voterPoll{}, errors.New("vr poll does not exist")
}

func (vr Voter) UpdatePoll(newvp voterPoll) (Voter, error) {
	pollid := newvp.PollID
	for i, vp := range vr.VoteHistory {
		if vp.PollID == pollid {
			vr.VoteHistory[i] = newvp
			return vr, nil
		}
	}
	return Voter{}, errors.New("vr poll does not exist")
}

func (vr Voter) DeletePoll(pollid uint) (Voter, error) {
	for i, vp := range vr.VoteHistory {
		if vp.PollID == pollid {
			vr.VoteHistory = append(vr.VoteHistory[:i], vr.VoteHistory[i+1:]...)
			return vr, nil
		}
	}
	return Voter{}, errors.New("vr poll does not exist")
}

func (vr Voter) ToJson() VoterJson {
	vj := VoterJson{
		Voter:     "/voters/" + strconv.FormatUint(uint64(vr.VoterID), 10),
		FirstName: vr.FirstName,
		LastName:  vr.LastName,
		History:   make([]string, 0),
	}
	pollPrefix := vj.Voter + "/polls/"
	for _, p := range vr.VoteHistory {
		vj.History = append(vj.History, pollPrefix+strconv.FormatUint(uint64(p.PollID), 10))
	}
	return vj
}

func (vr Voter) GetPollJson(vp voterPoll) voterPollJson {
	return voterPollJson{
		VoterPoll: vr.ToJson().Voter + "/polls/" + strconv.FormatUint(uint64(vp.PollID), 10),
		Date:      vp.VoteDate.String(),
	}
}
