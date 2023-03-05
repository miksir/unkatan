package deploy

import (
	"container/list"
	"encoding/json"
)

type deployHistory struct {
	list *list.List
	deep int
}

type historySaveStruct struct {
	CommandName string
	Command     json.RawMessage
}

func initDeployHistory(deep int) *deployHistory {
	d := deployHistory{
		deep: deep,
	}
	d.list = list.New()
	return &d
}

func (h *deployHistory) pushCommand(item DeployCommandI) {
	if h.list.Len() >= h.deep {
		e := h.list.Front()
		h.list.Remove(e)
	}
	h.list.PushBack(item)
}

func (h *deployHistory) GetList() []DeployCommandI {
	historySlice := make([]DeployCommandI, 0, h.list.Len())
	for e := h.list.Front(); e != nil; e = e.Next() {
		history := e.Value.(DeployCommandI)
		historySlice = append(historySlice, history)
	}
	return historySlice
}

func (h *deployHistory) saveState() ([]byte, error) {
	hList := h.GetList()
	saveList := make([]historySaveStruct, 0, len(hList))
	for _, item := range hList {
		comm, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		saveList = append(saveList, historySaveStruct{
			CommandName: CommandName(item),
			Command:     comm,
		})
	}
	jsonB, err := json.Marshal(saveList)
	return jsonB, err
}

func (h *deployHistory) restoreState(data []byte) error {
	historySlice := make([]historySaveStruct, 0, h.deep)
	err := json.Unmarshal(data, &historySlice)
	if err != nil {
		return err
	}
	h.list.Init()
	for _, item := range historySlice {
		var command DeployCommandI
		if item.CommandName == DeployActionOn {
			command = &DeployOnCommand{}
		} else if item.CommandName == DeployActionOff {
			command = &DeployOffCommand{}
		}
		err := json.Unmarshal(item.Command, command)
		if err != nil {
			return err
		}
		h.pushCommand(command)
	}
	return nil
}
