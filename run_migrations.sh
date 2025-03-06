#!/bin/bash
migrate -path db/migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/kpidb" up