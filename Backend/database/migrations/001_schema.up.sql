-- Department Table
CREATE TABLE deptData (
    deptID SERIAL PRIMARY KEY,
    deptName VARCHAR(255) NOT NULL UNIQUE
);

INSERT INTO deptData (deptName) VALUES ('CSE');

-- Admin Table
CREATE TABLE adminData (
    adminID SERIAL PRIMARY KEY,
    emailID VARCHAR(255) UNIQUE NOT NULL,
    adminName VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO adminData (emailID, adminName)
VALUES ('tharunkumarra@gmail.com', 'Tharun Kumar');
INSERT INTO adminData (emailID, adminName)
VALUES ('naganathan1555@gmail.com', 'Naganathan M R');

--Section table
CREATE TABLE sectionData (
    sectionID SERIAL PRIMARY KEY,
    sectionName VARCHAR(10) UNIQUE NOT NULL
);

INSERT INTO sectionData (sectionName) VALUES ('A');
INSERT INTO sectionData (sectionName) VALUES ('B');
INSERT INTO sectionData (sectionName) VALUES ('C');
INSERT INTO sectionData (sectionName) VALUES ('D');
INSERT INTO sectionData (sectionName) VALUES ('E');
INSERT INTO sectionData (sectionName) VALUES ('F');

--Semester table
CREATE TABLE semesterData (
    semesterID SERIAL PRIMARY KEY,
    semesterNumber INT UNIQUE NOT NULL
);

INSERT INTO semesterData (semesterNumber) VALUES (1);
INSERT INTO semesterData (semesterNumber) VALUES (2);
INSERT INTO semesterData (semesterNumber) VALUES (3);
INSERT INTO semesterData (semesterNumber) VALUES (4);
INSERT INTO semesterData (semesterNumber) VALUES (5);
INSERT INTO semesterData (semesterNumber) VALUES (6);
INSERT INTO semesterData (semesterNumber) VALUES (7);
INSERT INTO semesterData (semesterNumber) VALUES (8);

-- Faculty Table
CREATE TABLE facultyData (
    facultyID SERIAL PRIMARY KEY,
    emailID VARCHAR(255) UNIQUE NOT NULL,
    facultyName VARCHAR(255) NOT NULL,
    deptID INT REFERENCES deptData(deptID) ON DELETE SET NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO facultyData (emailID, facultyName, deptID)
VALUES ('naganathan1555@gmail.com', 'Naganathan', 1);

-- Student Table
CREATE TABLE studentData (
    studentID SERIAL PRIMARY KEY,
    rollNumber VARCHAR(50) UNIQUE NOT NULL,
    emailID VARCHAR(255) UNIQUE NOT NULL,
    studentName VARCHAR(255) NOT NULL,
    startYear INT NOT NULL,
    endYear INT NOT NULL,
    deptID INT REFERENCES deptData(deptID) ON DELETE SET NULL,
    sectionID INT REFERENCES sectionData(sectionID) ON DELETE SET NULL,
    semesterID INT REFERENCES semesterData(semesterID) ON DELETE SET NULL
);

INSERT INTO studentData (rollNumber, emailID, studentName, startYear, endYear, deptID, sectionID, semesterID)
VALUES ('CB.EN.U4CSE22240', 'naganathan1555@gmail.com', 'Naganathan', 2022, 2026, 1, 1, 1);
INSERT INTO studentData (rollNumber, emailID, studentName, startYear, endYear, deptID, sectionID, semesterID)
VALUES ('CB.EN.U4CSE22253', 'tharunkumarra@gmail.com', 'Tharun Kumar', 2022, 2026, 1, 1, 1);
INSERT INTO studentData (rollNumber, emailID, studentName, startYear, endYear, deptID, sectionID, semesterID)
VALUES ('CB.EN.U4CSE22222', 'hariprasathm777@gmail.com', 'Hari Prasath', 2022, 2026, 1, 1, 1);
INSERT INTO studentData (rollNumber, emailID, studentName, startYear, endYear, deptID, sectionID, semesterID)
VALUES ('CB.EN.U4CSE22245', 'bharathshan16@gmail.com', 'Bharath', 2022, 2026, 1, 1, 1);

-- Course Table
CREATE TABLE courseData (
    courseID SERIAL PRIMARY KEY,
    courseCode VARCHAR(50) UNIQUE NOT NULL,
    courseName VARCHAR(255) NOT NULL,
    courseDeptID INT REFERENCES deptData(deptID) ON DELETE CASCADE,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    courseType VARCHAR(1) NOT NULL,
    updatedBy INT REFERENCES adminData(adminID) ON DELETE SET NULL
);

INSERT INTO courseData (courseCode, courseName, courseDeptID, courseType, updatedBy)
VALUES ('19CSE312', 'POPL', 1, '1', 1);

-- Course-Faculty Mapping Table
CREATE TABLE courseFaculty (
    classroomID SERIAL PRIMARY KEY,
    courseID INT REFERENCES courseData(courseID) ON DELETE CASCADE,
    facultyID INT REFERENCES facultyData(facultyID) ON DELETE CASCADE,
    sectionID INT REFERENCES sectionData(sectionID) ON DELETE SET NULL,
    semesterID INT REFERENCES semesterData(semesterID) ON DELETE SET NULL,
    deptID INT REFERENCES deptData(deptID) ON DELETE SET NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    createdBy INT REFERENCES adminData(adminID) ON DELETE SET NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedBy INT REFERENCES adminData(adminID) ON DELETE SET NULL
);

INSERT INTO courseFaculty (courseID, facultyID, sectionID, semesterID, deptID, createdBy, updatedBy)
VALUES (1, 1, 3, 6, 1, 1, 1);

-- Meeting Table
CREATE TABLE meetingData (
    meetingID SERIAL PRIMARY KEY,
    startTime TIMESTAMP NOT NULL,
    endTime TIMESTAMP NOT NULL,
    classroomID INT REFERENCES courseFaculty(classroomID) ON DELETE CASCADE,
    meetingLink VARCHAR(255),
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    createdBy INT REFERENCES facultyData(facultyID) ON DELETE SET NULL,
    meetingDescription TEXT
);

INSERT INTO meetingData (startTime, endTime, classroomID, meetingLink, createdBy, meetingDescription) VALUES 
('2025-03-10 09:00:00', '2025-03-10 10:30:00', 1, 'https://meet.example.com/abc123', 1, 'Project kickoff meeting for spring semester');

-- Attendance Table
-- ('0' - Absent, '1' - Present)
CREATE TABLE attendanceData (
    meetingID INT REFERENCES meetingData(meetingID) ON DELETE CASCADE,
    studentID INT REFERENCES studentData(studentID) ON DELETE CASCADE,
    isPresent VARCHAR(1) DEFAULT '0',
    PRIMARY KEY (meetingID, studentID)
);

-- Quiz Table
CREATE TABLE quizData (
    quizID SERIAL PRIMARY KEY,
    classroomID INT REFERENCES courseFaculty(classroomID) ON DELETE CASCADE,
    quizName VARCHAR(255) NOT NULL,
    quizDescription TEXT,
    quizData TEXT,
    isOpenForAll VARCHAR(1) DEFAULT '0',
    startTime TIMESTAMP NOT NULL,
    endTime TIMESTAMP NOT NULL,
    quizDuration INT NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    createdBy INT REFERENCES facultyData(facultyID) ON DELETE SET NULL
);

-- Quiz Submission Table
CREATE TABLE quizSubmissionData (
    quizSubmissionID SERIAL PRIMARY KEY,
    quizID INT REFERENCES quizData(quizID) ON DELETE CASCADE,
    quizSubmissionData TEXT,
    studentID INT REFERENCES studentData(studentID) ON DELETE CASCADE,
    marks INT,
    submissionTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    isPublished VARCHAR(1) DEFAULT '0',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Assignment Table
CREATE TABLE assignmentData (
    assignmentID SERIAL PRIMARY KEY,
    classroomID INT REFERENCES courseFaculty(classroomID) ON DELETE CASCADE,
    assignmentName VARCHAR(255) NOT NULL,
    assignmentDescription TEXT,
    assignmentData TEXT,
    startTime TIMESTAMP NOT NULL,
    endTime TIMESTAMP NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedBy INT REFERENCES facultyData(facultyID) ON DELETE SET NULL
);

-- Assignment Submission Table
CREATE TABLE assignmentSubmission (
    assignmentSubmissionID SERIAL PRIMARY KEY,
    assignmentID INT REFERENCES assignmentData(assignmentID) ON DELETE CASCADE,
    studentID INT REFERENCES studentData(studentID) ON DELETE CASCADE,
    marks INT,
    isPublished BOOLEAN DEFAULT FALSE,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    createdBy INT REFERENCES facultyData(facultyID) ON DELETE SET NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Faculty Sample Data
INSERT INTO facultyData (emailID, facultyName, deptID) VALUES
('faculty1@college.com', 'Dr. Alice Johnson', 1),
('faculty2@college.com', 'Dr. Bob Smith', 1);

-- Course Sample Data
INSERT INTO courseData (courseCode, courseName, courseDeptID, courseType, updatedBy) VALUES
('CS101', 'Data Structures', 1, 'C', 1),
('CS102', 'Operating Systems', 1, 'C', 1);

-- Course-Faculty Mapping Sample Data
INSERT INTO courseFaculty (courseID, facultyID, sectionID, semesterID, deptID, createdBy, updatedBy) VALUES
(1, 1, 1, 1, 1, 1, 1),
(2, 2, 1, 1, 1, 1, 1);

INSERT INTO quizData (classroomID, quizName, quizDescription, quizData, isOpenForAll, startTime, endTime, quizDuration, createdBy) VALUES
(1, 'Quiz 1', 'Data Structures Basics', '{"Q1":"What is a stack?","Q2":"Explain linked list"}', '1', '2024-03-10 10:00:00', '2024-03-10 11:00:00', 60, 1),
(2, 'Quiz 2', 'OS Fundamentals', '{"Q1":"What is a process?","Q2":"Explain threading"}', '0', '2024-03-15 14:00:00', '2024-03-15 15:00:00', 60, 2);


-- Faculty Sample Data (Adding 15 More)
INSERT INTO facultyData (emailID, facultyName, deptID) VALUES
('tharunkumarra@gmail.com', 'Dr. Daniel White', 1),
('hariprasathm777@gmail.com', 'Dr. Eve Black', 1),
('bharathshan16@gmail.com', 'Dr. Frank Green', 1),
('faculty7@college.com', 'Dr. Grace Hall', 1),
('faculty8@college.com', 'Dr. Hannah King', 1),
('faculty9@college.com', 'Dr. Ian Lewis', 1),
('faculty10@college.com', 'Dr. Jack Martin', 1),
('faculty11@college.com', 'Dr. Kelly Nelson', 1),
('faculty12@college.com', 'Dr. Liam O’Connor', 1),
('faculty13@college.com', 'Dr. Mia Patel', 1),
('faculty14@college.com', 'Dr. Noah Quinn', 1),
('faculty15@college.com', 'Dr. Olivia Roberts', 1),
('faculty16@college.com', 'Dr. Peter Scott', 1),
('faculty17@college.com', 'Dr. Quinn Taylor', 1);

-- Course Sample Data (Adding 15 More)
INSERT INTO courseData (courseCode, courseName, courseDeptID, courseType, updatedBy) VALUES
('CS103', 'Computer Networks', 1, 'C', 1),
('CS104', 'Database Systems', 1, 'C', 1),
('CS105', 'Software Engineering', 1, 'C', 1),
('CS106', 'Web Technologies', 1, 'C', 1),
('CS107', 'Machine Learning', 1, 'C', 1),
('CS108', 'Artificial Intelligence', 1, 'C', 1),
('CS109', 'Cyber Security', 1, 'C', 1),
('CS110', 'Mobile Computing', 1, 'C', 1),
('CS111', 'Internet of Things', 1, 'C', 1),
('CS112', 'Blockchain Technology', 1, 'C', 1),
('CS113', 'Cloud Computing', 1, 'C', 1),
('CS114', 'Computer Vision', 1, 'C', 1),
('CS115', 'Bioinformatics', 1, 'C', 1),
('CS116', 'Big Data Analytics', 1, 'C', 1),
('CS117', 'Embedded Systems', 1, 'C', 1);

-- Course-Faculty Mapping Sample Data (Adding 15 More)
INSERT INTO courseFaculty (courseID, facultyID, sectionID, semesterID, deptID, createdBy, updatedBy) VALUES
(3, 3, 1, 1, 1, 1, 1),
(4, 4, 1, 1, 1, 1, 1),
(5, 5, 1, 1, 1, 1, 1),
(6, 6, 1, 1, 1, 1, 1),
(7, 7, 1, 1, 1, 1, 1),
(8, 8, 1, 1, 1, 1, 1),
(9, 9, 1, 1, 1, 1, 1),
(10, 10, 1, 1, 1, 1, 1),
(11, 11, 1, 1, 1, 1, 1),
(12, 12, 1, 1, 1, 1, 1),
(13, 13, 1, 1, 1, 1, 1),
(14, 14, 1, 1, 1, 1, 1),
(15, 15, 1, 1, 1, 1, 1),
(16, 16, 1, 1, 1, 1, 1),
(17, 17, 1, 1, 1, 1, 1);

-- Quiz Sample Data (Adding 15 More)
INSERT INTO quizData (classroomID, quizName, quizDescription, quizData, isOpenForAll, startTime, endTime, quizDuration, createdBy) VALUES
(3, 'Computer Networks Basics', 'A quiz covering fundamentals of networking.', '{"questions": ["What is an IP address?", "Explain OSI model layers."]}', '1', '2024-03-05 10:00:00', '2024-03-05 11:00:00', 60, 3),
(4, 'Database Queries', 'Quiz on SQL queries and normalization.', '{"questions": ["What is a primary key?", "Explain different SQL joins."]}', '0', '2024-03-06 14:00:00', '2024-03-06 15:00:00', 60, 4),
(5, 'Software Engineering Principles', 'A quiz covering SDLC models and agile methodology.', '{"questions": ["What are the phases of SDLC?", "Explain Scrum framework."]}', '1', '2024-03-07 10:00:00', '2024-03-07 11:00:00', 45, 5),
(6, 'HTML & CSS Basics', 'A beginner-level quiz on HTML elements and CSS properties.', '{"questions": ["What is a div element?", "Difference between ID and Class selectors."]}', '1', '2024-03-08 14:00:00', '2024-03-08 15:00:00', 30, 6),
(7, 'Introduction to Machine Learning', 'Covers ML algorithms and applications.', '{"questions": ["What is supervised learning?", "Explain K-means clustering."]}', '0', '2024-03-09 10:00:00', '2024-03-09 11:00:00', 45, 7),
(8, 'Deep Learning Concepts', 'Quiz covering neural networks and deep learning models.', '{"questions": ["What is a perceptron?", "Explain backpropagation."]}', '0', '2024-03-10 14:00:00', '2024-03-10 15:00:00', 50, 8),
(9, 'Cyber Security Fundamentals', 'Covers cryptography, network security, and attacks.', '{"questions": ["What is symmetric encryption?", "Explain phishing attack."]}', '1', '2024-03-11 10:00:00', '2024-03-11 11:00:00', 45, 9),
(10, 'Mobile Computing Overview', 'Basic quiz on mobile computing and wireless networks.', '{"questions": ["What is mobile IP?", "Explain Bluetooth technology."]}', '1', '2024-03-12 14:00:00', '2024-03-12 15:00:00', 40, 10),
(11, 'Internet of Things (IoT)', 'Quiz covering IoT devices and protocols.', '{"questions": ["What is MQTT?", "Explain IoT architecture."]}', '0', '2024-03-13 10:00:00', '2024-03-13 11:00:00', 35, 11),
(12, 'Blockchain Fundamentals', 'Covers decentralization and blockchain consensus mechanisms.', '{"questions": ["What is a smart contract?", "Explain Proof of Work."]}', '1', '2024-03-14 14:00:00', '2024-03-14 15:00:00', 50, 12),
(13, 'Cloud Computing Basics', 'Quiz covering cloud service models.', '{"questions": ["What is SaaS?", "Explain difference between IaaS and PaaS."]}', '0', '2024-03-15 10:00:00', '2024-03-15 11:00:00', 40, 13),
(14, 'Computer Vision Techniques', 'Covers image processing and AI-based vision techniques.', '{"questions": ["What is convolution in CNN?", "Explain edge detection."]}', '1', '2024-03-16 14:00:00', '2024-03-16 15:00:00', 45, 14),
(15, 'Big Data Analytics', 'A quiz covering Hadoop, Spark, and data processing.', '{"questions": ["What is MapReduce?", "Explain role of Apache Spark."]}', '0', '2024-03-17 10:00:00', '2024-03-17 11:00:00', 55, 15),
(16, 'Bioinformatics Applications', 'Quiz covering computational biology.', '{"questions": ["What is BLAST in bioinformatics?", "Explain genomics and proteomics."]}', '1', '2024-03-18 14:00:00', '2024-03-18 15:00:00', 50, 16),
(17, 'Embedded Systems Basics', 'Quiz on microcontrollers and real-time OS.', '{"questions": ["What is an RTOS?", "Explain embedded Linux."]}', '0', '2024-03-19 10:00:00', '2024-03-19 11:00:00', 35, 17);