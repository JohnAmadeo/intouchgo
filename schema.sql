DROP TABLE letters;
DROP TABLE inmates;
DROP TABLE facilities;

-- Assumes that drafts are stored on-device and never reach the cloud
CREATE TABLE facilities (
    name VARCHAR PRIMARY KEY,
    shortName VARCHAR,
    address VARCHAR NOT NULL CHECK (length(address) > 0)
);

CREATE TABLE inmates (
    id VARCHAR UNIQUE,
    state VARCHAR,
    inmateNumber VARCHAR,
    firstName VARCHAR NOT NULL CHECK (length(firstName) > 0),
    lastName VARCHAR NOT NULL CHECK (length(lastName) > 0),
    dateOfBirth DATE,
    facility VARCHAR REFERENCES facilities(name),
    active BOOLEAN NOT NULL,
    PRIMARY KEY(state, inmateNumber)
);

CREATE TABLE letters (
    id VARCHAR PRIMARY KEY,
    author VARCHAR NOT NULL CHECK (length(author) > 0),
    recipient VARCHAR NOT NULL REFERENCES inmates(id),
    subject VARCHAR,
    text VARCHAR NOT NULL CHECK (length(text) > 0),
    timeSent DATE NOT NULL,
    timeLastEdited DATE NOT NULL,
    isDraft BOOLEAN NOT NULL
);

INSERT INTO facilities VALUES
    ('Bridgeport Correctional Center', 'Bridgeport CC', '1106 North Avenue, Bridgeport, CT 06604'),
    ('Brooklyn Correctional Institution', 'Brooklyn CI', '59 Hartford Road, Brooklyn, CT 06234'),
    ('Cheshire Correctional Institution', 'Cheshire CI', '900 Highland Avenue, Cheshire, CT 06410 '),
    ('Corrigan-Radgowski Correctional Center', 'Corrigan-Radgowski CC', '986 Norwich-New London Turnpike, Uncasville, CT 06382 '),
    ('Garner Correctional Institution', 'Garner CI', '50 Nunnawauk Road, Newtown, CT 06470 '),
    ('Hartford Correctional Center', 'Hartford CC', '177 Weston Street, Hartford, CT 06120 '),
    ('MacDougall-Walker Correctional Institution', 'MacDougall-Walker CI', '1153 East Street, South Suffield , CT 06080 '),
    ('Manson Youth Institution', 'Manson YI', '42 Jarvis Street, Cheshire, CT 06410'),
    ('New Haven Correctional Center', 'New Haven CC', '245 Whalley Avenue, POB 8000, New Haven, CT 06511 '),
    ('Northern Correctional Institution', 'Northern CI', '287 Bilton Road, POB 665, Somers, CT 06071'),
    ('Osborn Correctional Institution', 'Osborn CI', '335 Bilton Road, POB 100, Somers, CT 06071'),
    ('Robinson Correctional Institution', 'Robinson CI', '285 Shaker Road, POB 1400, Enfield , CT 06082'),
    ('Willard-Cybulski Correctional Institution / Cybulski Community Reintegration Center', 'Willard-Cybulski CI', '391 Shaker Road, Enfield , CT 06082'),
    ('York Correctional Institution', 'York CI', '201 West Main Street, Niantic, CT 06357');