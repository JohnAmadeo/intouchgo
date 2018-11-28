DROP TABLE letters;
DROP TABLE inmates;
DROP TABLE facilities;

-- Assumes that drafts are stored on-device and never reach the cloud
CREATE TABLE facilities (
    name VARCHAR PRIMARY KEY,
    shortName VARCHAR,
    addressLine VARCHAR NOT NULL CHECK (length(addressLine) > 0),
    city VARCHAR NOT NULL CHECK(length(city) > 0),
    state VARCHAR NOT NULL CHECK(length(state) > 0),
    zip VARCHAR NOT NULL CHECK(length(zip) > 0),
    lobTestAddressId VARCHAR NOT NULL CHECK(length(lobTestAddressId) > 0),
    lobLiveAddressId VARCHAR NOT NULL -- CHECK(length(lobLiveAddressId) > 0)
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
    timeDeliveredEstimate DATE NOT NULL,
    isDraft BOOLEAN NOT NULL,
    lobLetterId VARCHAR NOT NULL CHECK (length(lobLetterId) > 0)
);

INSERT INTO facilities VALUES
    ('Bridgeport Correctional Center',                                                      'Bridgeport CC',            '1106 North Avenue',                'Bridgeport',       'CT',   '06604', 'adr_8ee3fb884ac685d4', ''),
    ('Brooklyn Correctional Institution',                                                   'Brooklyn CI',              '59 Hartford Road',                 'Brooklyn',         'CT',   '06234', 'adr_55a1e5bb48f6a073', ''),
    ('Cheshire Correctional Institution',                                                   'Cheshire CI',              '900 Highland Avenue',              'Cheshire',         'CT',   '06410', 'adr_00f18daf8fde8d69', ''),
    ('Corrigan-Radgowski Correctional Center',                                              'Corrigan-Radgowski CC',    '986 Norwich-New London Turnpike',  'Uncasville',       'CT',   '06382', 'adr_cc3b40e138076a1a', ''),
    ('Garner Correctional Institution',                                                     'Garner CI',                '50 Nunnawauk Road',                'Newtown',          'CT',   '06470', 'adr_3b2be4a08bb48254', ''),
    ('Hartford Correctional Center',                                                        'Hartford CC',              '177 Weston Street',                'Hartford',         'CT',   '06120', 'adr_248557789145db55', ''),
    ('MacDougall-Walker Correctional Institution',                                          'MacDougall-Walker CI',     '1153 East Street',                 'South Suffield',   'CT',   '06080', 'adr_35e76be2e0ae8e33', ''),
    ('Manson Youth Institution',                                                            'Manson YI',                '42 Jarvis Street',                 'Cheshire',         'CT',   '06410', 'adr_e4831e813280d467', ''),
    ('New Haven Correctional Center',                                                       'New Haven CC',             '245 Whalley Avenue, POB 8000',     'New Haven',        'CT',   '06511', 'adr_9f2f31176ec1e8fa', 'adr_5e9e8a215c4ef466'),
    ('Northern Correctional Institution',                                                   'Northern CI',              '287 Bilton Road, POB 665',         'Somers',           'CT',   '06071', 'adr_bca8ef16fa248624', ''),
    ('Osborn Correctional Institution',                                                     'Osborn CI',                '335 Bilton Road, POB 100',         'Somers',           'CT',   '06071', 'adr_cdf0735da6bbc41b', ''),
    ('Robinson Correctional Institution',                                                   'Robinson CI',              '285 Shaker Road, POB 1400',        'Enfield',          'CT',   '06082', 'adr_58ef03fe53326182', ''),
    ('Willard-Cybulski Correctional Institution / Cybulski Community Reintegration Center', 'Willard-Cybulski CI',      '391 Shaker Road',                  'Enfield',          'CT',   '06082', 'adr_abeaafa9bf1f31eb', ''),
    ('York Correctional Institution',                                                       'York CI',                  '201 West Main Street',             'Niantic',          'CT',   '06357', 'adr_88ad7d8c101fced4', '');
    
    