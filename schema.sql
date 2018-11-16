DROP TABLE facilities;
DROP TABLE inmates;
DROP TABLE letters;

-- Assumes that drafts are stored on-device and never reach the cloud
CREATE TABLE facilities (
    name VARCHAR PRIMARY KEY,
    address VARCHAR NOT NULL CHECK (length(text) > 0)
);

CREATE TABLE inmates (
    id VARCHAR PRIMARY KEY,
    firstName VARCHAR NOT NULL CHECK (length(text) > 0),
    lastName VARCHAR NOT NULL CHECK (length(text) > 0),
    inmateNumber VARCHAR NOT NULL CHECK (length(text) > 0),
    dateOfBirth DATE,
    facility VARCHAR REFERENCES facilities(name)
);

CREATE TABLE letters (
    id VARCHAR PRIMARY KEY,
    author VARCHAR NOT NULL CHECK (length(text) > 0),
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
    
INSERT INTO inmates
    VALUES ('asdf-123s-ddss', 'John', 'Grant', 'AA123', '05/03/91', 'CT Prison 1');
    
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