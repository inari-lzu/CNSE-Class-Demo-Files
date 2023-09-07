package poll

import (
	"errors"
	"strconv"
)

type pollOption struct {
	PollOptionID    uint   `json:"id"`
	PollOptionValue string `json:"value"`
}

type pollOptionJson struct {
	Option string `json:"option"`
	Value  string `json:"value"`
}

type Poll struct {
	PollID       uint         `json:"id"`
	PollTitle    string       `json:"title"`
	PollQuestion string       `json:"question"`
	PollOptions  []pollOption `json:"options,omitempty"`
}

type PollJson struct {
	Poll     string   `json:"poll"`
	Title    string   `json:"title"`
	Question string   `json:"question"`
	Options  []string `json:"options,omitempty"`
}

func (p Poll) GetID() uint {
	return p.PollID
}

func NewPoll() Poll {
	p := Poll{}
	p.PollOptions = make([]pollOption, 0)
	return p
}

func (p Poll) Update(pl Poll) (Poll, error) {
	pl.PollID = p.PollID // PollOptions is also can be changed in updating
	return pl, nil
}

func (p Poll) UpdateOptions(pl Poll) (Poll, error) {
	p.PollOptions = pl.PollOptions
	return p, nil
}

func (p Poll) GetAllOptions() []pollOption {
	return p.PollOptions
}

func (p Poll) DeleteAllOptions() Poll {
	p.PollOptions = make([]pollOption, 0)
	return p
}

func (p Poll) AddPoll(newpo pollOption) (Poll, error) {
	for _, po := range p.PollOptions {
		if newpo.PollOptionID == po.PollOptionID {
			return Poll{}, errors.New("pollOption already exists")
		}
	}
	p.PollOptions = append(p.PollOptions, newpo)
	return p, nil
}

func (p Poll) GetOption(optionId uint) (pollOption, error) {
	for _, po := range p.PollOptions {
		if po.PollOptionID == optionId {
			return po, nil
		}
	}
	return pollOption{}, errors.New("pollOption does not exist")
}

func (p Poll) UpdateOption(newpo pollOption) (Poll, error) {
	optionId := newpo.PollOptionID
	for i, po := range p.PollOptions {
		if po.PollOptionID == optionId {
			p.PollOptions[i] = newpo
			return p, nil
		}
	}
	return Poll{}, errors.New("pollOption does not exist")
}

func (p Poll) DeleteOption(optionId uint) (Poll, error) {
	for i, po := range p.PollOptions {
		if po.PollOptionID == optionId {
			p.PollOptions = append(p.PollOptions[:i], p.PollOptions[i+1:]...)
			return p, nil
		}
	}
	return Poll{}, errors.New("pollOption does not exist")
}

func (p Poll) ToJson() PollJson {
	pj := PollJson{
		Poll:     "/polls/" + strconv.FormatUint(uint64(p.PollID), 10),
		Title:    p.PollTitle,
		Question: p.PollQuestion,
		Options:  make([]string, 0),
	}
	optionPrefix := pj.Poll + "/options/"
	for _, op := range p.PollOptions {
		pj.Options = append(pj.Options, optionPrefix+strconv.FormatUint(uint64(op.PollOptionID), 10))
	}
	return pj
}

func (p Poll) GetOptionJson(po pollOption) pollOptionJson {
	return pollOptionJson{
		Option: p.ToJson().Poll + "/options/" + strconv.FormatUint(uint64(po.PollOptionID), 10),
		Value:  po.PollOptionValue,
	}
}
