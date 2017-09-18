CREATE TABLE customer (
	id sequence 
	name text
	environmenttype text // test | prod
);

CREATE TABLE logsmeta (
	id integer
	customer int
	identitfier text
	label text
	guid text
	product_version text
	install_id int
	change_type text
	releases text
	created_at timestamp without time zone
	updated_at timestamp without time zone
	failed_product text
	failed_release text
);


CREATE TABLE logs (
	id sequence
	logtext bytea 
);


releases: {[
	Name: string
	Version: string
]}