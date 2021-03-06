create table "matches"
(
    id          integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    ip          VARCHAR(15) NOT NULL,
    started_at  bigint      NOT NULL,
    map         VARCHAR(50) NOT NULL,
    rounds      smallint    NOT NULL,
    duration    integer     NOT NULL,
    won         bool        NOT NULL default false,
    inserted_at bigint      NOT NULL DEFAULT date_part('epoch'::text, now()),
    UNIQUE (ip, started_at, map)
)

create table "users"
(
    id               bigint PRIMARY KEY,
    name             VARCHAR(32) NOT NULL,
    avatar_hash      CHAR(40)             DEFAULT NULL,
    kills            integer     NOT NULL default 0,
    deaths           integer     NOT NULL default 0,
    fratricide       integer     NOT NULL default 0,
    kd               numeric(10, 2)       DEFAULT NULL,
    all_weapon_stats jsonb       NOT NULL default '{}'::jsonb,
    inserted_at      bigint      NOT NULL DEFAULT date_part('epoch'::text, now())
)

CREATE INDEX idx_users_kills
    ON users (kills);

create table "match_user_stats"
(
    match_id     integer NOT NULL,
    user_id      bigint  NOT NULL,
    kills        integer NOT NULL default 0,
    deaths       integer NOT NULL default 0,
    fratricide   integer NOT NULL default 0,
    weapon_stats jsonb   NOT NULL default '{}'::jsonb,

    FOREIGN KEY (match_id) REFERENCES matches (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    UNIQUE (match_id, user_id)
)

create table "user_medals"
(
    user_id     bigint  NOT NULL,
    medal_id    integer NOT NULL,
    value       integer          default NULL,
    current     bool    NOT NULL default false,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    UNIQUE (user_id, medal_id)
)


update users
set kills = a.total from (select user_id, sum(kills) as total from match_user_stats group by user_id) a
WHERE users.id = a.user_id;

update users
set deaths = a.total from (select user_id, sum(deaths) as total from match_user_stats group by user_id) a
WHERE users.id = a.user_id;

update users
set fratricide = a.total from (select user_id, sum(fratricide) as total from match_user_stats group by user_id) a
WHERE users.id = a.user_id;

update users
set kd = cast(kills as decimal) / deaths
where kills > 100
  and deaths != 0;

update users
set kd = 9999
where kills > 100
  and deaths = 0

update users
set all_weapon_stats = stats.agg from (
select user_id, jsonb_object_agg(k, val) as agg
from (
         select user_id, k, sum(v::numeric) as val
         from match_user_stats
                  join lateral jsonb_each_text(weapon_stats) j(k, v) on true
         group by user_id, k
     ) tt
group by user_id) stats
where user_id = id
