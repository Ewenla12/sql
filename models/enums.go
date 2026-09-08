package models

// StudentStatus mirrors the Postgres ENUM `student_status`.
type StudentStatus string

const (
    StatusActive     StudentStatus = "active"
    StatusGraduated  StudentStatus = "graduated"
    StatusSuspended  StudentStatus = "suspended"
    StatusDroppedOut StudentStatus = "dropped_out"
    StatusExpelled   StudentStatus = "expelled"
)

// StudentLevel mirrors the Postgres ENUM `student_level`.
type StudentLevel string

const (
    Level100 StudentLevel = "100"
    Level200 StudentLevel = "200"
    Level300 StudentLevel = "300"
    Level400 StudentLevel = "400"
    Level500 StudentLevel = "500"
)

// Semester mirrors the Postgres ENUM `semester`.
type Semester string

const (
    SemesterFirst  Semester = "first"
    SemesterSecond Semester = "second"
    SemesterSummer Semester = "summer"
)

// AdminRole mirrors the Postgres ENUM `admin_role`.
type AdminRole string

const (
    RoleHOD            AdminRole = "hod"
    RoleFacultyOfficer AdminRole = "faculty_officer"
    RoleDean           AdminRole = "dean"
    RoleRegistrar      AdminRole = "registrar"
    RoleSuperAdmin     AdminRole = "super_admin"
)

// ApplicationStatus mirrors the Postgres ENUM `application_status`.
type ApplicationStatus string

const (
    ApplicationPending  ApplicationStatus = "pending"
    ApplicationApproved ApplicationStatus = "approved"
    ApplicationRejected ApplicationStatus = "rejected"
)
