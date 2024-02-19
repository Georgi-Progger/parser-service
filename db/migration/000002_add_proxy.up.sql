CREATE TABLE IF NOT EXISTS "proxies"(
    id SERIAL PRIMARY KEY,
    body VARCHAR(255),
    isactive BOOLEAN NOT NULL 
);

INSERT INTO proxies (body, isactive) VALUES ('194.8.232.46:4153',true);
INSERT INTO proxies (body, isactive) VALUES ('198.41.206.65:80',true);
INSERT INTO proxies (body, isactive) VALUES ('104.18.247.214:80',true);
INSERT INTO proxies (body, isactive) VALUES ('185.238.228.201:80',true);
INSERT INTO proxies (body, isactive) VALUES ('172.67.54.200:80',true);
