TRUNCATE TABLE course_registration, student_courses, lecturers, students, admin, courses, departments, faculties
RESTART IDENTITY CASCADE;

INSERT INTO faculties (faculty_name) VALUES
('Faculty of Science'), ('Faculty of Engineering'), ('Faculty of Arts'), ('Faculty of Social Sciences');

INSERT INTO departments (faculty_id, department_name, department_code) VALUES
(1, 'Computer Science', 'CSC'), (1, 'Mathematics', 'MTH'), (1, 'Physics', 'PHY'), (1, 'Chemistry', 'CHM'),
(2, 'Electrical Engineering', 'EEE'), (2, 'Mechanical Engineering', 'MEE'), (2, 'Civil Engineering', 'CVE'),
(3, 'English Language', 'ENG'), (3, 'History', 'HIS'), (4, 'Economics', 'ECO');

INSERT INTO admin (username, email, password_hash, role, department_id) VALUES
('jane.hod', 'jane.hod@school.edu', '$2a$10$vD6a4jYytFhGO309e6uYp.OJ.Ti72mMx8OOxktK4ccj5NKlQlHiqi', 'hod', 1),
('super.admin', 'admin@school.edu', '$2a$10$vD6a4jYytFhGO309e6uYp.OJ.Ti72mMx8OOxktK4ccj5NKlQlHiqi', 'super_admin', NULL);

DO $$
DECLARE
    dept RECORD;
    yr INT;
    years INT[] := ARRAY[2021, 2022, 2023, 2024];
    year_count INT;
    i INT;
    first_names TEXT[] := ARRAY['John','Mary','Chidi','Aisha','Femi','Blessing','Yusuf','Ngozi','Tunde','Zainab',
                                  'Emeka','Halima','Kunle','Grace','Ibrahim','Funmi','Segun','Amaka','Bala','Ronke',
                                  'Dele','Chinwe','Sani','Nkechi','Wale','Rita','Ken','Titi','Musa','Ifeoma',
                                  'Adamu','Yemisi','Chibuike','Esther','Kola','Maryam','Obinna','Sarah','David','Grace'];
    last_names TEXT[] := ARRAY['Balogun','Chukwu','Adekunle','Bello','Nwosu','Eze','Ibrahim','Okoro','Obi','Sule',
                                 'Ade','Udo','Musa','Ojo','Alade','Nnaji','Bakare','Aliyu','Onwu','Etim',
                                 'Yusuf','Fashola','Ogundipe','Umeh','Abdullahi','Obiora','Adebayo','Nwachukwu','Saro','Bankole',
                                 'Garba','Anya','Lawal','Ajayi','Okafor','Danjuma','Awolowo','Yakubu','Johnson','Okonkwo'];
    genders TEXT[] := ARRAY['male', 'female'];
    statuses TEXT[] := ARRAY['active','active','active','active','active','active','active','active',
                              'graduated','suspended','dropped_out','expelled'];
    fname TEXT;
    lname TEXT;
    matric TEXT;
    admission_yr INT;
    reference_year CONSTANT INT := 2025;
    yrs_in INT;
BEGIN
    FOR dept IN SELECT department_id, department_code FROM departments LOOP
        FOREACH yr IN ARRAY years LOOP
            year_count := 10 + floor(random() * 6)::INT;
            admission_yr := yr;
            yrs_in := reference_year - admission_yr;
            yrs_in := LEAST(GREATEST(yrs_in, 1), 4);

            FOR i IN 1..year_count LOOP
                fname := first_names[1 + floor(random() * array_length(first_names, 1))::INT];
                lname := last_names[1 + floor(random() * array_length(last_names, 1))::INT];
                matric := dept.department_code || '/' || admission_yr || '/' || LPAD(i::TEXT, 3, '0');

                INSERT INTO students (
                    students_id, first_name, last_name, department_id, level, status,
                    admission_year, gender, email, password_hash, cgpa
                ) VALUES (
                    matric, fname, lname, dept.department_id,
                    (yrs_in * 100)::TEXT::student_level,
                    statuses[1 + floor(random() * array_length(statuses, 1))::INT]::student_status,
                    admission_yr,
                    genders[1 + floor(random() * array_length(genders, 1))::INT],
                    lower(replace(matric, '/', '.')) || '@school.edu',
                    '$2a$10$vD6a4jYytFhGO309e6uYp.OJ.Ti72mMx8OOxktK4ccj5NKlQlHiqi',
                    ROUND((random() * 4)::NUMERIC, 2)
                );
            END LOOP;
        END LOOP;
    END LOOP;
END $$;