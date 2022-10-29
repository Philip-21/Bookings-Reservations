create_table("users") {
  t.Column("id","integer",{primary:true})
  t.Column("first_name", "string", {"default":""})
  t.Column("last_name", "string", {"default":""})
  t.Column("email", "string", {})
  t.Column("password", "string", {"size": 60})
  t.Column("access_level", "integer", {"default": 1})
  
}

create_table("reservation"){
    t.Column("id","integer",{primary:true})
    t.Column("first_name", "string", {"default":""})
    t.Column("last_name", "string", {"default":""})
    t.Column("email", "string", {})
    t.Column("phone","string",{"default":""})
    t.Column("start_date","date",{})
    t.Column("end_date","date",{})
    t.Column("room_id","integer",{})
}

create_table("rooms") {
	t.Column("id","integer", {primary:true})
	t.Column("room_name", "string", {"default": ""})
}

create_table("restrictions") {
	t.Column("id", "integer", {primary: true})
	t.Column("restriction_name", "string", {"default": ""})	
}

add_foreign_key("reservation","room_id",{"rooms":["id"]},{
    "on_delete":"cascade",
    "on_update":"cascade",
})

add_index("room_restrictions",["start_date","end_date"],{})
add_index("room_restrictions","room_id",{})
add_index("room_restrictions","reservation_id",{})
add_index("users","email",{"unique":true})
add_index("reservation", "email", {})
add_index("reservation", "last_name", {})

change_column("room_restrictions", "reservation_id", "integer", {"null":true})
add_column("reservation", "processed", "integer", {"default": 0})
drop_foreign_key("room_restrictions", "room_restrictions_restrictions_id_fk", {})
drop_table("restrictions")