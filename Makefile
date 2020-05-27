.PHONY: run
run:
	docker-compose run init_db
	docker-compose up agent resultsprocessor

.PHONY: reset
reset:
	docker-compose run reset_db


.PHONY: tests
tests:
	docker-compose run unit_tests


.PHONY: e2e_tests
e2e_tests:
	docker-compose run e2e_tests