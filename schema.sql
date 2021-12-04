create table "matches"
(
    id          integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    ip          VARCHAR(15) NOT NULL,
    started_at  bigint   NOT NULL,
    map         VARCHAR(50) NOT NULL,
    rounds      smallint    NOT NULL,
    duration    integer     NOT NULL,
    won         bool        NOT NULL default false,
    inserted_at bigint   NOT NULL DEFAULT date_part('epoch'::text, now()),
    UNIQUE (ip, started_at, map)
)

create table "users"
(
    id               bigint PRIMARY KEY,
    kills            integer   NOT NULL default 0,
    deaths           integer   NOT NULL default 0,
    kd               numeric(10, 2) DEFAULT NULL,
    all_weapon_stats jsonb     NOT NULL default '{}'::jsonb,
    inserted_at      bigint   NOT NULL DEFAULT date_part('epoch'::text, now())
)

create table "match_user_stats"
(
    match_id     integer NOT NULL,
    user_id      bigint  NOT NULL,
    kills        integer NOT NULL default 0,
    deaths       integer NOT NULL default 0,
    weapon_stats jsonb   NOT NULL default '{}'::jsonb,

    FOREIGN KEY (match_id) REFERENCES matches (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    UNIQUE (match_id, user_id)
)


update users
set kills = a.total
    from (select user_id, sum(kills) as total from match_user_stats group by user_id) a
WHERE users.id = a.user_id;

update users
set deaths = a.total
    from (select user_id, sum(deaths) as total from match_user_stats group by user_id) a
WHERE users.id = a.user_id;

update users
set kd = cast(kills as decimal)/deaths
where kills > 100 and deaths != 0;

update users
set kd = 9999
where kills > 100 and deaths = 0

