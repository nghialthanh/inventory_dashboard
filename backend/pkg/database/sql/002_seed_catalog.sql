INSERT INTO makes (id, name) VALUES
    (1, 'Ford'),
    (2, 'Honda'),
    (3, 'Toyota')
ON CONFLICT (id) DO NOTHING;

INSERT INTO models (id, make_id, name) VALUES
    (1, 1, 'F-150'),
    (2, 1, 'Mustang'),
    (3, 2, 'Accord'),
    (4, 2, 'Civic'),
    (5, 3, 'Camry'),
    (6, 3, 'Corolla')
ON CONFLICT (id) DO NOTHING;

SELECT setval(pg_get_serial_sequence('makes', 'id'), (SELECT MAX(id) FROM makes));
SELECT setval(pg_get_serial_sequence('models', 'id'), (SELECT MAX(id) FROM models));
