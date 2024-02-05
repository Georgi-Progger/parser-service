CREATE TABLE IF NOT EXISTS "annoucement"(
    id SERIAL PRIMARY KEY,
    model VARCHAR(255),
    price VARCHAR(255),
    year VARCHAR(255),
    generation VARCHAR(255),
    mileage VARCHAR(255),
    history VARCHAR(255),
    pts VARCHAR(255),
    owners VARCHAR(255),
    condition VARCHAR(255),
    modification VARCHAR(255),
    engine_volume VARCHAR(255), 
    engine_type VARCHAR(255),
    transmission VARCHAR(255),
    drive VARCHAR(255),
    equipment VARCHAR(255),
    body_type VARCHAR(255),
    color VARCHAR(255),
    steering VARCHAR(255),
    vin VARCHAR(255),
    exchange VARCHAR(255)
);