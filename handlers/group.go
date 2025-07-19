package handlers

import (
	"net/http"
	"telcohub/db"
	"telcohub/models"
	"telcohub/utils"
	"text/template"
)

func ShowGroups(w http.ResponseWriter, r *http.Request) {
	user, _ := utils.GetUserFromSession(r, store)

	var groups []models.Group
	db.DB.Preload("Owner").Find(&groups)

	data := []map[string]interface{}{}
	haveGroup := false
	createGroup := true
	for _, g := range groups {
		var members []models.GroupUser
		db.DB.Preload("User").Where("group_id = ?", g.ID).Find(&members)

		memberList := []map[string]interface{}{}
		for _, m := range members {
			memberList = append(memberList, map[string]interface{}{
				"Username": m.User.Username,
				"ID":       m.ID,
				"IsAdmin":  m.IsAdmin,
			})
		}
		haveGroup = g.OwnerID == user.ID
		if haveGroup {
			createGroup = false
		}
		data = append(data, map[string]interface{}{
			"ID":           g.ID,
			"Name":         g.Name,
			"OwnerName":    g.Owner.Username,
			"Members":      memberList,
			"IsGroupOwner": haveGroup,
		})
	}
	tmpl := template.Must(template.ParseFiles("templates/groups.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Groups":      data,
		"CreateGroup": createGroup,
	})
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	user, _ := utils.GetUserFromSession(r, store)
	name := r.FormValue("name")

	group := models.Group{Name: name, OwnerID: user.ID}
	db.DB.Create(&group)

	db.DB.Create(&models.GroupUser{
		UserID:  user.ID,
		GroupID: group.ID,
		IsAdmin: true,
	})

	http.Redirect(w, r, "/groups/list", http.StatusSeeOther)
}

func AddToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := utils.GetParamID(r)
	username := r.FormValue("username")

	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	db.DB.Create(&models.GroupUser{
		UserID:  user.ID,
		GroupID: groupID,
		IsAdmin: false,
	})

	http.Redirect(w, r, "/groups/list", http.StatusSeeOther)
}

func ToggleAdmin(w http.ResponseWriter, r *http.Request) {
	user, _ := utils.GetUserFromSession(r, store)
	memberID := utils.GetParamID(r)

	var member models.GroupUser
	if err := db.DB.First(&member, memberID).Error; err != nil {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	// 🛡️ Only the group owner can toggle admin
	var group models.Group
	if err := db.DB.First(&group, member.GroupID).Error; err != nil || group.OwnerID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	member.IsAdmin = !member.IsAdmin
	db.DB.Save(&member)

	w.WriteHeader(http.StatusOK)
}

func RemoveMember(w http.ResponseWriter, r *http.Request) {
	user, _ := utils.GetUserFromSession(r, store)
	memberID := utils.GetParamID(r)

	var member models.GroupUser
	if err := db.DB.First(&member, memberID).Error; err != nil {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	// 🛡️ Only group owner can remove members (except self)
	var group models.Group
	if err := db.DB.First(&group, member.GroupID).Error; err != nil || group.OwnerID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	db.DB.Delete(&member)
	w.WriteHeader(http.StatusOK)
}
