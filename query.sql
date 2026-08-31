CREATE TYPE student_status AS ENUM ('active', 'graduated', 'suspended', 'dropped_out', 'expelled');
CREATE TYPE student_level AS ENUM ('100', '200', '300', '400', '500');
CREATE TYPE semester AS ENUM ('first', 'second', 'summer');
CREATE TYPE admin_role AS ENUM ('hod', 'faculty_officer', 'dean', 'registrar', 'super_admin');
CREATE TYPE application_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE faculties (
    faculty_id SERIAL PRIMARY KEY,
    faculty_name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE departments (
    department_id SERIAL PRIMARY KEY,
    faculty_id INT NOT NULL,
    department_name VARCHAR(100) NOT NULL UNIQUE,
    department_code VARCHAR(10) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (faculty_id) REFERENCES faculties(faculty_id)
);

CREATE TABLE courses (
    course_id SERIAL PRIMARY KEY,
    department_id INT NOT NULL,  
    course_code VARCHAR(10) NOT NULL UNIQUE,
    course_title VARCHAR(100) NOT NULL,
    credit_units INT NOT NULL DEFAULT 3,
    level INT NOT NULL DEFAULT 100,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (department_id) REFERENCES departments(department_id)
);

CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    students_id VARCHAR(20) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE,
    department_id INT NOT NULL,
    level student_level NOT NULL DEFAULT '100',
    status student_status NOT NULL DEFAULT 'active',
    admission_year INT NOT NULL,
    gender VARCHAR(10),
    email VARCHAR(100) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    cgpa DECIMAL(5, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (department_id) REFERENCES departments(department_id)
);

CREATE TABLE student_courses (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(20) NOT NULL,
    course_id INT NOT NULL,
    academic_session VARCHAR(10) NOT NULL,
    score DECIMAL(5, 2),
    grade VARCHAR(2),
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (student_id) REFERENCES students(students_id),
    FOREIGN KEY (course_id) REFERENCES courses(course_id),
    UNIQUE (student_id, course_id, academic_session)
);

CREATE TABLE admin(
    admin_id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role admin_role NOT NULL DEFAULT 'hod',
    department_id INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (department_id) REFERENCES departments(department_id)
);

CREATE TABLE course_registration(
    registration_id SERIAL PRIMARY KEY,
    student_id VARCHAR(20) NOT NULL,
    course_id INT NOT NULL,
    academic_session VARCHAR(10) NOT NULL,
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewed_by INT NULL,
    status application_status NOT NULL DEFAULT 'pending',

    FOREIGN KEY (student_id) REFERENCES students(students_id),
    FOREIGN KEY (course_id) REFERENCES courses(course_id),
    FOREIGN KEY (reviewed_by) REFERENCES admin(admin_id),
    UNIQUE (student_id, course_id, academic_session)
);

CREATE TABLE lecturers (
    lecturer_id SERIAL PRIMARY KEY,
    staff_id VARCHAR(20) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    department_id INT NOT NULL,
    email VARCHAR(100) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (department_id) REFERENCES departments(department_id)
);

