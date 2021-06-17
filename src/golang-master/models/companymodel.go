package models

import (
	"database/sql"
	"log"
	"github.com/jmoiron/sqlx"
	"golang-master/lang"
	"fmt"
)

type Company struct {
    Id	   int64 `json:"id"`
    Name   string `json:"name"`
	Status int64 `json:"status"`
}


type Companies []Company

type DataTablesCompanies struct {
	RecordsTotal int64 `json:"recordsTotal"`
	Companies Companies `json:"companies"`
}

type ReqCompany struct {
    Name   string `json:"name" validate:"required,min=2,max=100,alpha_space"`
	Status int64 `json:"status" validate:"required"`
}

// GetCompanies returns companies list
func GetCompanies(db *sql.DB) *Companies {
	companies := Companies{}

	res, err := db.Query("SELECT id,name,status FROM companies")

    defer res.Close()

    if err != nil {
        log.Fatal(err)
    }
	for res.Next() {
		var company Company
        err := res.Scan(&company.Id, &company.Name, &company.Status)
		if err != nil {
            log.Fatal(err)
        }
		companies =  append(companies, company)
	
    }
	return &companies

}

func GetCompaniesSqlxDataTables(db *sqlx.DB, reqcompany *DataTablesRequest) *DataTablesCompanies {
	dataTablesFields := map[int]string { 0:"id", 1:"name",2:"status" }
	datatablecompanies := DataTablesCompanies{}
	orderByQuery := ``;	
	for _,service := range reqcompany.Order {
			fmt.Println(service.Column)
            if len(orderByQuery) > 0 {
				orderByQuery =  orderByQuery + `, `+dataTablesFields[service.Column] + ` ` + service.Dir;
			} else {
				orderByQuery = ` Order By `+dataTablesFields[service.Column] + ` ` + service.Dir;
			}	
    }
	fmt.Println(orderByQuery)
	searchQuery := ``;	
	if(reqcompany.Search.Value != "") {
		searchQuery = ` AND name like '%`+reqcompany.Search.Value+`%'`;
	}	
	err := db.Select(&datatablecompanies.Companies, "SELECT SQL_CALC_FOUND_ROWS id,name,status FROM companies WHERE 1=1 "+searchQuery+" "+orderByQuery+" LIMIT ? OFFSET ? ", reqcompany.Length,reqcompany.Start)
	err = db.Get(&datatablecompanies.RecordsTotal,"SELECT count(0) FROM companies WHERE  1=1 "+searchQuery+"")

	if err != nil {
		log.Fatal(err)
	}
	return &datatablecompanies

}

func GetCompaniesSqlx(db *sqlx.DB) *Companies {
	companies := Companies{}

	err := db.Select(&companies, "SELECT id,name,status FROM companies")

    if err != nil {
        log.Fatal(err)
    }
	return &companies

}

// PostCompanySqlx insert company
func PostCompanySqlx(db *sqlx.DB, reqcompany *ReqCompany) (*Company, string) {
	name := reqcompany.Name
	status := reqcompany.Status
	
	var company Company
	
	stmt, err := db.Prepare("INSERT INTO companies(name,status) VALUES(?,?)")
	if err != nil {
		return &company, ErrHandler(err)
	}
	result, err := stmt.Exec(name,status)
	if err != nil {
		return &company, ErrHandler(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return &company, ErrHandler(err)
	}
	
	err=db.Get(&company, "SELECT id,name,status FROM companies where id = ?",id)
	if err != nil {
		return &company, lang.Get("no_result")
	}
	return &company, ""
}

// GetCompany get company
func GetCompany(db *sqlx.DB, id string) (*Company, string) {
	var company Company
	err:=db.Get(&company, "SELECT id,name,status FROM companies where id = ?",id)
	if err != nil {
		return &company, lang.Get("no_result")
	}
	return &company, ""
}


// PostCompanySqlx insert company
func EditCompany(db *sqlx.DB, reqcompany *ReqCompany,id int64) (*Company, string) {
	name := reqcompany.Name
	status := reqcompany.Status
	
	var company Company

	stmt, err := db.Prepare("Update companies set name=?, status=? where id = ?")
	if err != nil {
		return &company, ErrHandler(err)
	}
	_, err = stmt.Exec(name,status,id)
	if err != nil {
		return &company, ErrHandler(err)
	}

	err=db.Get(&company, "SELECT id,name,status FROM companies where id = ?",id)
	if err != nil {
		return &company, lang.Get("no_result")
	}
	return &company, ""
}

// DeleteCompany get company
func DeleteCompany(db *sqlx.DB, id string) string {
	stmt, err := db.Prepare("DELETE FROM companies where id = ?")
	if err != nil {
		return ErrHandler(err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return ErrHandler(err)
	}
	return ""
}


