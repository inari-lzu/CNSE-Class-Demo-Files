package vote

import (
	"strconv"
	"time"
)

type Vote struct {
	VoteID    uint `json:"id"`
	VoterID   uint `json:"voterId"`
	PollID    uint `json:"pollId"`
	VoteValue uint `json:"choiceId"`
}

type Links struct {
	Vote      string `json:"vote"`
	Voter     string `json:"voter"`
	VoterPoll string `json:"voterPoll"`
	Poll      string `json:"poll"`
	Choice    string `json:"choice"`
}

func (v Vote) GetID() uint {
	return v.VoteID
}

func NewVote() Vote {
	v := Vote{}
	return v
}

func (v Vote) Update(newv Vote) (Vote, error) {
	v.VoteValue = newv.VoteValue
	return v, nil
}

func (v *Vote) ToLinks(hostName string, voterUrl string, pollUrl string) Links {
	return Links{
		Vote:      hostName + "/votes/" + strconv.FormatUint(uint64(v.VoteID), 10),
		Voter:     voterUrl + "/voters/" + strconv.FormatUint(uint64(v.VoterID), 10),
		VoterPoll: voterUrl + "/voters/" + strconv.FormatUint(uint64(v.VoterID), 10) + "/polls/" + strconv.FormatUint(uint64(v.PollID), 10),
		Poll:      pollUrl + "/polls/" + strconv.FormatUint(uint64(v.PollID), 10),
		Choice:    pollUrl + "/polls/" + strconv.FormatUint(uint64(v.PollID), 10) + "/options/" + strconv.FormatUint(uint64(v.VoteValue), 10),
	}
}

func (v *Vote) ToVoteHistoryRecord() string {
	return `{"history": [{"id": ` + strconv.FormatUint(uint64(v.PollID), 10) +
		`, "date": "` + time.Now().Format("2006-01-02T15:04:05Z07:00") + `"}]}`
}
