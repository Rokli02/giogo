package server

import "giogo/ui/pages/minesweeper/model"

const (
	websocket_status_path = "/socket/status"
	websocket_action_path = "/socket"
	server_status         = "/status"
	server_health_check   = "/health-check"
)

const (
	count_of_port_reservation_tries = 5
	max_message_size_in_bytes       = 16_384
	max_size_of_mines_in_data       = max_message_size_in_bytes - max_message_size_in_bytes%model.SizeOfMineElementInBytes
)

const (
	Private_Host = "localhost"
	Public_Host  = "0.0.0.0"
)
