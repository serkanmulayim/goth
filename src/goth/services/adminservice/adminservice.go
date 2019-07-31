package adminservice

import (
	"context"
	"goth/objects/admin"

	"go.etcd.io/etcd/clientv3"
)

const (
	etcdAdminServicePasswordPath = "/admin/password/"
	etcdAdminServiceProfilePath  = "/admin/profile/"
	etcdAdminServiceAdminCounter = "/admin/count"
	adminCountLimit              = 50
)

func putPassword(email string, UserID int64, password string, cli *clientv3.Client) error {
	return nil
}

func putProfile(admin *admin.Object) error {

	return nil
}

//UpdateProfile updates an admin profile.
//returns user does not exist error message or internal server error messages
func UpdateProfile(admin *admin.Object, id int64, password string, cli *clientv3.Client) error {

	//find profile by ID from etcdAdminServiceProfile
	//if not found return user does not exist error

	//check if email is changed
	//if yes change the email in etcdAdminServicePassword table by deleting the old row etc.
	//if password is provided update it with the new password

	//update profile with the new profile in etcdAdminServiceProfile table

	return nil
}

//DeleteProfile deletes admin with id parameters
func DeleteProfile(id int64, requestorId int64, cli *clientv3.Client) error {
	//check if id and requestorid are the same. If yes, return user cannot delete itself error.
	//This would ensure that there would be at least 1 admin

	//delete admin from both tables
}

//CreateProfile creates an admin profile. There is a limit of "adminCountLimit"
//create profile at etcdAdminServiceProfilePath          (  id/ {profile})
//create username/password etcdAdminServicePasswordPath  (  email/{id, SHA256(pwd)}  )
func CreateProfile(admin *admin.Object, password string, cli *clientv3.Client) error {

	//get admin counter (this will be the userid)
	//increment admin counter
	//check if username exists in etcdAdminServicePasswordPath. if yes return user exists error
	//create profile at etcdAdminServiceProfilePath          (  id/ {profile})
	//create username/password etcdAdminServicePasswordPath  (  email/{id, SHA256(pwd)}  )

	return nil
}

//CreateAdminIfNone creates an admin account with email/username:"admin" password:password id=1
func CreateAdminIfNone(cli *clientv3.Client) error {
	//Scan admin profiles table.
	//if it is empty create a new admin with id:1 email:admin firstname:admin lastname:admin password:password
	ctx, cancel := context.WithTimeout(context.Background(), etc)
	return nil
}
