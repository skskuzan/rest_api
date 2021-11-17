package crud

import (
	"context"
	"fmt"
	"rest_api/database"
	"rest_api/grpc/pkg/api"
	"rest_api/vars"
	"strconv"
)

type GRPCServer struct{}

func (s *GRPCServer) AddStd(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	var std, std1 vars.Student
	std.Name = req.GetName()
	std.Age = req.GetAge()
	std.Mail = req.GetMail()

	msg := "New Student created:"

	rows := database.ProceedQ("SELECT * FROM people")
	max := 0
	for rows.Next() {
		rows.Scan(&std1.Id, &std1.Name, &std1.Mail, &std1.Age)
		mid, _ := strconv.Atoi(std1.Id)
		if mid > max {
			max = mid
		}
	}

	err := database.ExecuteQ("INSERT INTO people (id,name, mail,age) VALUES (?,?,?,?)", max+1, std.Name, std.Mail, std.Age)
	if err != nil {
		fmt.Println(err)
		msg = fmt.Sprint(err)
	}

	return &api.AddResponse{Result: msg + fmt.Sprint(std)}, nil
}

func (s *GRPCServer) ReadStd(ctx context.Context, req *api.ReadRequest) (*api.ReadResponse, error) {

	msg := "Here is that Student:"

	rows := database.ProceedQ("SELECT * FROM people WHERE id =?", fmt.Sprint(req.GetId()))

	rows.Next()
	p := vars.Student{}
	err := rows.Scan(&p.Id, &p.Name, &p.Mail, &p.Age)

	if err != nil {
		fmt.Println(err)
		msg = fmt.Sprint(err)
	}

	return &api.ReadResponse{Result: msg + fmt.Sprint(p)}, nil

}

func (s *GRPCServer) ChangeStd(ctx context.Context, req *api.ChangeRequest) (*api.ChangeResponse, error) {
	var std vars.Student
	std.Id = fmt.Sprint(req.GetId())
	std.Name = req.GetName()
	std.Age = req.GetAge()
	std.Mail = req.GetMail()

	msg := "Student " + std.Id + " changed: "

	err := database.ExecuteQ("UPDATE people SET name=?,mail=?,age=? WHERE id=?", std.Name, std.Mail, std.Age, std.Id)

	if err != nil {
		fmt.Println(err)
		msg = fmt.Sprint(err)
	}

	return &api.ChangeResponse{Result: msg + fmt.Sprint(std)}, nil

}

func (s *GRPCServer) DeleteStd(ctx context.Context, req *api.DelRequest) (*api.DelResponse, error) {
	err := database.ExecuteQ("DELETE FROM people WHERE id=" + fmt.Sprint(req.GetId()))
	msg := "Student " + fmt.Sprint(req.GetId()) + " was deleted."

	if err != nil {
		fmt.Println(err)
		msg = fmt.Sprint(err)
	}
	return &api.DelResponse{Result: msg}, nil

}

func (s *GRPCServer) ReadDB(ctx context.Context, req *api.DBRequest) (*api.DBResponse, error) {

	_, err := vars.Database.Query("select * from people")
	msg := "Database:"

	if err != nil {
		fmt.Println(err)
		msg = fmt.Sprint(err)
	}

	rows := database.ProceedQ("SELECT * FROM people")

	p1 := make([]vars.Student, 0)
	for rows.Next() {
		p2 := vars.Student{}
		err := rows.Scan(&p2.Id, &p2.Name, &p2.Mail, &p2.Age)
		if err != nil {
			fmt.Println(err)
			continue
		}
		p1 = append(p1, p2)
	}

	return &api.DBResponse{Result: msg + fmt.Sprint(p1)}, nil

}
