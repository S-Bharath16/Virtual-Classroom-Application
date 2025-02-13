CREATE TABLE IF NOT EXISTS userRole (
    roleId SERIAL PRIMARY KEY,
    roleName VARCHAR(255) NOT NULL
);

INSERT INTO userRole (roleName) VALUES ('ADMIN'), ('DEPT_HEAD'), ('OFFICE'), ('PROFESSOR'), ('Nags');

CREATE TABLE IF NOT EXISTS departmentData (
    deptId SERIAL PRIMARY KEY,
    deptName VARCHAR(255) NOT NULL
);

INSERT INTO departmentData (deptName) VALUES ('CSE'), ('AIE'), ('CYS');

CREATE TABLE IF NOT EXISTS managementData (
    managerId SERIAL PRIMARY KEY,
    managerEmail VARCHAR(255) NOT NULL UNIQUE,
    managerPassword VARCHAR(255) NOT NULL,
    deptId INT NOT NULL,
    roleId INT NOT NULL,
    managerFullName VARCHAR(255),
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    createdBy INT NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    managerStatus CHAR(1) DEFAULT '1',
    FOREIGN KEY (deptId) REFERENCES departmentData(deptId),
    FOREIGN KEY (roleId) REFERENCES userRole(roleId),
    FOREIGN KEY (createdBy) REFERENCES managementData(managerId)
);

CREATE TABLE IF NOT EXISTS studentData (
    studentId SERIAL PRIMARY KEY,
    studentName VARCHAR(255) NOT NULL,
    studentRollNumber VARCHAR(50) UNIQUE NOT NULL,
    studentGender CHAR(1) NOT NULL DEFAULT 'N',
    studentPhone CHAR(10) NOT NULL,
    studentEmail VARCHAR(255) UNIQUE NOT NULL,
    studentPassword VARCHAR(255) NOT NULL,
    studentDob CHAR(10) NOT NULL,
    studentDeptId INT NOT NULL,
    studentSection CHAR(1) NOT NULL,
    studentBatchStart CHAR(4) NOT NULL,
    studentBatchEnd CHAR(4) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    createdBy INT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    studentStatus CHAR(1) DEFAULT '1',
    FOREIGN KEY (studentDeptId) REFERENCES departmentData(deptId),
    FOREIGN KEY (createdBy) REFERENCES managementData(managerId)
);

CREATE TABLE IF NOT EXISTS forgotPasswordStudent (
    studentId INT NOT NULL,
    otp VARCHAR(6) NOT NULL,
    expiresAt TIMESTAMP DEFAULT NOW() + INTERVAL '5 minutes',
    FOREIGN KEY (studentId) REFERENCES studentData(studentId)
);

CREATE TABLE IF NOT EXISTS forgotPasswordManagement (
    managerId INT NOT NULL,
    otp VARCHAR(6) NOT NULL,
    expiresAt TIMESTAMP DEFAULT NOW() + INTERVAL '5 minutes',
    FOREIGN KEY (managerId) REFERENCES managementData(managerId)
);

CREATE TABLE IF NOT EXISTS courseData (
    courseId SERIAL PRIMARY KEY,
    courseCode VARCHAR(20) UNIQUE NOT NULL,
    courseName VARCHAR(255) NOT NULL,
    courseDeptId INT NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    createdBy INT NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    courseStatus CHAR(1) DEFAULT '1',
    courseType CHAR(1) DEFAULT '1',
    updatedBy INT NOT NULL,
    FOREIGN KEY (courseDeptId) REFERENCES departmentData(deptId),
    FOREIGN KEY (createdBy) REFERENCES managementData(managerId),
    FOREIGN KEY (updatedBy) REFERENCES managementData(managerId)
);

CREATE TABLE IF NOT EXISTS courseFaculty (
    classroomId SERIAL PRIMARY KEY,
    courseId INT NOT NULL,
    managerId INT NOT NULL,
    batchStart CHAR(4) NOT NULL,
    batchEnd CHAR(4) NOT NULL,
    section CHAR(1) NOT NULL,
    isMentor CHAR(1) DEFAULT '0',
    isActive CHAR(1) DEFAULT '1',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    createdBy INT NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedBy INT NOT NULL,
    FOREIGN KEY (courseId) REFERENCES courseData(courseId),
    FOREIGN KEY (managerId) REFERENCES managementData(managerId),
    FOREIGN KEY (createdBy) REFERENCES managementData(managerId),
    FOREIGN KEY (updatedBy) REFERENCES managementData(managerId)
);

CREATE TABLE IF NOT EXISTS studentCourse (
    studentId INT NOT NULL,
    classroomId INT NOT NULL,
    FOREIGN KEY (studentId) REFERENCES studentData(studentId),
    FOREIGN KEY (classroomId) REFERENCES courseFaculty(classroomId)
);

CREATE TABLE IF NOT EXISTS classRoomData (
    classId SERIAL PRIMARY KEY,
    classroomId INT NOT NULL,
    classStartTime TIMESTAMP NOT NULL,
    classEndTime TIMESTAMP NOT NULL,
    classLink VARCHAR(255) NOT NULL,
    classStatus CHAR(1) DEFAULT '1',
    FOREIGN KEY (classroomId) REFERENCES courseFaculty(classroomId)
);

CREATE TABLE IF NOT EXISTS attendanceData (
    classId INT NOT NULL,
    studentId INT NOT NULL,
    isPresent CHAR(1) DEFAULT '0',
    FOREIGN KEY (classId) REFERENCES classRoomData(classId),
    FOREIGN KEY (studentId) REFERENCES studentData(studentId)
);

CREATE TABLE IF NOT EXISTS quizData (
    quizId SERIAL PRIMARY KEY,
    classroomId INT NOT NULL,
    quizName VARCHAR(255) NOT NULL,
    quizDescription TEXT NOT NULL,
    quizData JSON NOT NULL,
    isOpenForAll CHAR(1) DEFAULT '0',
    startTime TIMESTAMP NOT NULL,
    endTime TIMESTAMP NOT NULL,
    duration VARCHAR(20),
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    createdBy INT NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedBy INT NOT NULL,
    FOREIGN KEY (classroomId) REFERENCES courseFaculty(classroomId),
    FOREIGN KEY (createdBy) REFERENCES managementData(managerId),
    FOREIGN KEY (updatedBy) REFERENCES managementData(managerId)
);

CREATE TABLE IF NOT EXISTS quizSubmission (
    quizSubmissionId SERIAL PRIMARY KEY,
    quizId INT NOT NULL,
    quizSubmissionData JSON NOT NULL,
    studentId INT NOT NULL,
    marks INT,
    isPublished CHAR(1) DEFAULT '0',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    createdBy INT NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedBy INT NOT NULL,
    FOREIGN KEY (quizId) REFERENCES quizData(quizId),
    FOREIGN KEY (studentId) REFERENCES studentData(studentId),
    FOREIGN KEY (createdBy) REFERENCES managementData(managerId),
    FOREIGN KEY (updatedBy) REFERENCES managementData(managerId)
);

CREATE TABLE IF NOT EXISTS assignmentData (
    assignmentId SERIAL PRIMARY KEY,
    classroomId INT NOT NULL,
    assignmentName VARCHAR(255) NOT NULL,
    assignmentDescription TEXT NOT NULL,
    isOpenToAll CHAR(1) DEFAULT '0',
    startTime TIMESTAMP NOT NULL,
    endTime TIMESTAMP NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    createdBy INT NOT NULL,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedBy INT NOT NULL,
    FOREIGN KEY (classroomId) REFERENCES courseFaculty(classroomId),
    FOREIGN KEY (createdBy) REFERENCES managementData(managerId),
    FOREIGN KEY (updatedBy) REFERENCES managementData(managerId)
);