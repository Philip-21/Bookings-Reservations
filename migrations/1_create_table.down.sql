drop_table("users")

drop_table("reservation")

drop_table("rooms")

drop_table("restrictions")

drop_foreign_key("reservation",reservation_room_id_fk)
drop_table("room_restrictions)


drop_foreign_key("room_restrictions", "room_restrictions_rooms_id_fk", {})
drop_foreign_key("room_restrictions", "room_restrictions_restrictions_id_fk", {})
drop_foreign_key("room_restrictions", "room_restrictions_reservation_id_fk", {})
drop_index("users","users_email_idx")

drop_index("reservation", "reservation_email_idx")
drop_index("reservation", "reservation_last_name_idx")
drop_column("reservation", "processed")
drop_column("room_restrictions","room_restrictions_restrictions_id",{})
drop_table("room_restrictions)


drop_foreign_key("room_restrictions", "room_restrictions_rooms_id_fk", {})
drop_foreign_key("room_restrictions", "room_restrictions_restrictions_id_fk", {})
drop_foreign_key("room_restrictions", "room_restrictions_reservation_id_fk", {})