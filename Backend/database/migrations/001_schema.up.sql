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
VALUES ('CB.EN.U4CSE22240', 'naganathan1555@gmail.com', 'Naganathan', 2022, 2026, 1, 3, 6);
INSERT INTO studentData (rollNumber, emailID, studentName, startYear, endYear, deptID, sectionID, semesterID)
VALUES ('CB.EN.U4CSE22253', 'tharunkumarra@gmail.com', 'Tharun Kumar', 2022, 2026, 1, 3, 6);
INSERT INTO studentData (rollNumber, emailID, studentName, startYear, endYear, deptID, sectionID, semesterID)
VALUES ('CB.EN.U4CSE22222', 'hariprasathm777@gmail.com', 'Hari Prasath', 2022, 2026, 1, 3, 6);
INSERT INTO studentData (rollNumber, emailID, studentName, startYear, endYear, deptID, sectionID, semesterID)
VALUES ('CB.EN.U4CSE22245', 'bharathshan16gmail.com', 'Bharath', 2022, 2026, 1, 3, 6);

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