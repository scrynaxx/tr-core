CREATE TABLE "order".vehicles
(
    id                  uuid,
    type                text           NOT NULL,
    registration_number text           NOT NULL,
    length_meters       numeric(10, 2) NOT NULL,
    width_meters        numeric(10, 2) NOT NULL,
    height_meters       numeric(10, 2) NOT NULL,
    capacity_tonnes     numeric(10, 2) NOT NULL,
    created_at          timestamptz    NOT NULL DEFAULT now(),
    updated_at          timestamptz    NOT NULL DEFAULT now(),
    archived_at         timestamptz,
    CONSTRAINT pk_vehicles PRIMARY KEY (id)
);

CREATE UNIQUE INDEX uidx_vehicles_registration_number ON "order".vehicles (registration_number) WHERE registration_number IS NOT NULL;


CREATE TABLE "order".offerings
(
    id          uuid,
    name        text           NOT NULL,
    price       numeric(12, 2) NOT NULL,
    modifiers   text[]         NOT NULL,
    created_at  timestamptz    NOT NULL DEFAULT now(),
    updated_at  timestamptz    NOT NULL DEFAULT now(),
    archived_at timestamptz,
    CONSTRAINT pk_offerings PRIMARY KEY (id)
);


CREATE TABLE "order".cargo_packages
(
    id          uuid,
    name        text        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    archived_at timestamptz,
    CONSTRAINT pk_cargo_packages PRIMARY KEY (id)
);


CREATE TABLE "order".packing_materials
(
    id          uuid,
    name        text           NOT NULL,
    price       numeric(12, 2) NOT NULL,
    created_at  timestamptz    NOT NULL DEFAULT now(),
    updated_at  timestamptz    NOT NULL DEFAULT now(),
    archived_at timestamptz,
    CONSTRAINT pk_packing_materials PRIMARY KEY (id)
);


INSERT INTO "order".offerings (id, name, price, modifiers)
VALUES (gen_random_uuid(), 'Грузчики', 0, '{quantity}'),
       (gen_random_uuid(), 'Сборщики', 0, '{quantity}'),
       (gen_random_uuid(), 'Аренда авто с экипажем - S', 1000, '{}'),
       (gen_random_uuid(), 'Аренда авто с экипажем - M', 2000, '{}'),
       (gen_random_uuid(), 'Аренда авто с экипажем - L', 3000, '{}'),
       (gen_random_uuid(), 'Аренда авто с экипажем - XL', 4000, '{}'),
       (gen_random_uuid(), 'Аренда авто с экипажем - XXL', 5000, '{}');


INSERT INTO "order".packing_materials (id, name, price)
VALUES (gen_random_uuid(), 'Стретч пленка', 0),
       (gen_random_uuid(), 'Воздушно-пузырчатая плёнка', 0),
       (gen_random_uuid(), 'Лист гофро-картона', 0),
       (gen_random_uuid(), 'Коробки', 0),
       (gen_random_uuid(), 'Крафтовая бумага', 0),
       (gen_random_uuid(), 'Малярный скотч', 0),
       (gen_random_uuid(), 'Маркер', 0);


INSERT INTO "order".cargo_packages (id, name)
VALUES (gen_random_uuid(), 'Палет'),
       (gen_random_uuid(), 'Россыпь'),
       (gen_random_uuid(), 'Ящик');


INSERT INTO "order".vehicles (id, type, registration_number, length_meters, width_meters, height_meters,
                              capacity_tonnes)
VALUES (gen_random_uuid(), 'awning', 'А123ВС777', 3.10, 2.00, 2.10, 2.00),
       (gen_random_uuid(), 'awning', 'В456ЕК799', 4.20, 2.10, 2.20, 3.50),
       (gen_random_uuid(), 'fridge', 'К789МН197', 5.20, 2.40, 2.50, 5.00),
       (gen_random_uuid(), 'fridge', 'О321РТ777', 6.50, 2.50, 2.60, 10.00),
       (gen_random_uuid(), 'thermos', 'С654УХ799', 7.20, 2.50, 2.70, 15.00);