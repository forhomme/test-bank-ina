package adapters

const (
	sql_InsertUser = `INSERT INTO users(email, name, password) VALUES (?,?,?)`

	sql_ListUser = `SELECT id, email, name FROM users ORDER BY updated_at DESC`

	sql_UserById = `SELECT id, email, name FROM users WHERE id = ?`

	sql_GetUserByEmail = `SELECT id,email,name,password FROM users WHERE email =?`

	sql_UpdateUserById = `UPDATE users set email=?, name=?, password=? WHERE id=?`

	sql_DeleteUserById = `CALL sp_DeleteUser(?)`

	sql_InsertTask = `INSERT INTO tasks(user_id,title,description,status) VALUES (?,?,?,?)`

	sql_ListTask = `SELECT id,user_id,title,description,status FROM tasks ORDER BY updated_at DESC`

	sql_GetTaskById = `SELECT id,user_id,title,description,status FROM tasks WHERE id = ?`

	sql_UpdateTaskById = `UPDATE tasks set user_id=?,title=?,description=?,status=? WHERE id =?`

	sql_DeleteTaskById = `DELETE FROM tasks WHERE id = ?`
)
