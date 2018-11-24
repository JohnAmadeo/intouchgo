DROP TABLE letters;
DROP TABLE inmates;
DROP TABLE facilities;

-- Assumes that drafts are stored on-device and never reach the cloud
CREATE TABLE facilities (
    name VARCHAR PRIMARY KEY,
    address VARCHAR NOT NULL CHECK (length(address) > 0)
);

CREATE TABLE inmates (
    id VARCHAR UNIQUE,
    state VARCHAR,
    inmateNumber VARCHAR,
    firstName VARCHAR NOT NULL CHECK (length(firstName) > 0),
    lastName VARCHAR NOT NULL CHECK (length(lastName) > 0),
    dateOfBirth DATE,
    -- facility VARCHAR REFERENCES facilities(name),
    facility VARCHAR,
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

INSERT INTO facilities
    VALUES ('CT Prison 1', '5445 Hardaway Park Drive');
INSERT INTO facilities 
    VALUES ('CT Prison 2', '33A Holmes St NW');
INSERT INTO facilities 
    VALUES ('CT Prison 3', '7F One Lawrence Way');
    
INSERT INTO inmates VALUES
    ('asdf-123s-ddss', 'CT', 'AA123', 'John', 'Grant', '05/03/91', 'CT Prison 1', true),
    ('123a-1das-mmji', 'CT', 'AA296', 'Marlene', 'Enrique', '01/05/78', 'CT Prison 2', false);
    
INSERT INTO letters
    VALUES (
        'sd8f-1239-zxqw', 
        'jadk157', 
        'asdf-123s-ddss', 
        'Merry Christmas!', 
        'Hi John!',
        '02/10/18',
        '02/10/18',
        false 
    );