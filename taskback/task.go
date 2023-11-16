package taskback

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/iufansh/iufans/utils"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type TaskBack struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Total  int    `json:"total"`
	Done   int    `json:"done"`
	Remark string `json:"remark"`
}

func AddTaskBack(name string, total int) (int64, error) {
	task := TaskBack{
		Id:     time.Now().Unix(),
		Name:   name,
		Total:  total,
		Done:   0,
		Remark: "",
	}
	var taskKeys string
	if err := utils.GetCache("task_back_list", &taskKeys); err != nil {
		logs.Error("AddTask GetCache task_back_list err:", err)
		return 0, err
	}
	taskKey := genTaskKey(task.Id)
	taskKeys = taskKeys + "," + taskKey
	if err := utils.SetCache("task_back_list", taskKeys, 86400); err != nil {
		logs.Error("AddTask SetCache task_back_list err:", err)
		return 0, err
	}
	//b, _ := json.Marshal(task)
	if err := utils.SetCache(taskKey, task, 86400); err != nil {
		logs.Error("AddTask SetCache taskKey=", taskKey, " err:", err)
		return 0, err
	}

	return task.Id, nil
}

func genTaskKey(id int64) string {
	return fmt.Sprintf("task_back_%d", id)
}

func UpdateTaskBackStatus(id int64, total int, increment int, remark string) error {
	var taskKeys string
	if err := utils.GetCache("task_back_list", &taskKeys); err != nil {
		logs.Error("UpdateTaskStatus GetCache task_back_list err:", err)
		return err
	}
	taskKey := genTaskKey(id)
	if !strings.Contains(taskKeys, taskKey) {
		return errors.New("Task Id not exist")
	}
	var task TaskBack
	if err := utils.GetCache(taskKey, &task); err != nil {
		logs.Error("UpdateTaskStatus GetCache taskKey=", taskKey, " err:", err)
		return err
	}
	if total > 0 {
		task.Total = total
	}
	task.Done = task.Done + increment
	task.Remark = remark
	if err := utils.SetCache(taskKey, task, 86400); err != nil {
		logs.Error("UpdateTaskStatus SetCache taskKey=", taskKey, " err:", err)
		return err
	}
	return nil
}

func GetAllTaskBack() ([]TaskBack, error) {
	var list = make([]TaskBack, 0)
	var taskKeys string
	if err := utils.GetCache("task_back_list", &taskKeys); err != nil {
		logs.Error("UpdateTaskStatus GetCache task_back_list err:", err)
		return nil, err
	}
	if taskKeys == "" {
		return list, nil
	}
	keys := strings.Split(taskKeys, ",")
	for _, v := range keys {
		if v == "" {
			continue
		}
		var task TaskBack
		if err := utils.GetCache(v, &task); err != nil {
			logs.Error("UpdateTaskStatus GetCache taskKey=", v, " err:", err)
			task.Remark = "查询异常"
		}
		list = append(list, task)
	}
	return list, nil
}
