CREATE TABLE movie;
ALTER TABLE movie ADD COLUMN id integer;
ALTER TABLE movie ADD COLUMN star integer;
ALTER TABLE movie ADD COLUMN budget integer;
ALTER TABLE movie ADD COLUMN profit integer;
ALTER TABLE movie ADD COLUMN producer integer;
ALTER TABLE movie ADD COLUMN dialogue integer;

CREATE TABLE actor;
ALTER TABLE actor ADD COLUMN id integer;
ALTER TABLE actor ADD COLUMN name text;
ALTER TABLE actor ADD COLUMN native_lang integer;

CREATE TABLE producer;
ALTER TABLE producer ADD COLUMN id integer;
ALTER TABLE producer ADD COLUMN name text;
ALTER TABLE producer ADD COLUMN native_lang integer;

CREATE TABLE language;
ALTER TABLE language ADD COLUMN id integer;
ALTER TABLE language ADD COLUMN name text;

ALTER TABLE movie ADD CONSTRAINT fk_movie_actor FOREIGN KEY (star) REFERENCES actor (id)
ALTER TABLE movie ADD CONSTRAINT fk_movie_producer FOREIGN KEY (producer) REFERENCES producer (id)
ALTER TABLE movie ADD CONSTRAINT fk_movie_language FOREIGN KEY (dialogue) REFERENCES language (id)
ALTER TABLE producer ADD CONSTRAINT fk_producer_language FOREIGN KEY (native_lang) REFERENCES language (id)
ALTER TABLE actor ADD CONSTRAINT fk_actor_language FOREIGN KEY (native_lang) REFERENCES language (id)
