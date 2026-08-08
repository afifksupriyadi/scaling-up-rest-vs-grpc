package seeder

import (
	"math/rand"

	"scaling-up-rest-vs-grpc/internal/data/model"
)

const shapeLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// shapeFakeStr returns a random string of exactly n characters. Content is
// intentionally non-semantic (unlike Student's gofakeit-based data), since
// this experiment measures serialization cost by field count and length,
// not by data realism.
func shapeFakeStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = shapeLetters[rand.Intn(len(shapeLetters))]
	}
	return string(b)
}

func shapeFakeInt32(max int32) int32 { return int32(rand.Intn(int(max)) + 1) }
func shapeFakeInt64(max int64) int64 { return rand.Int63n(max) + 1 }
func shapeFakeFloat64() float64      { return rand.Float64() * 100 }
func shapeFakeBool() bool            { return rand.Intn(2) == 0 }

const (
	depth1WideEmployeesPerCompany = 40

	depth3NarrowDepartments      = 3
	depth3NarrowTeamsPerDept     = 2
	depth3NarrowEmployeesPerTeam = 4

	depth4WideDepartments      = 2
	depth4WideTeamsPerDept     = 2
	depth4WideEmployeesPerTeam = 3
	depth4WideCertsPerEmployee = 15
)

// ---------- Depth 0 ----------

func shapeFakeEmployeeFlat() *model.EmployeeFlat {
	return &model.EmployeeFlat{
		EmployeeId: shapeFakeStr(10), FullName: shapeFakeStr(18), FirstName: shapeFakeStr(10), LastName: shapeFakeStr(10),
		Gender: shapeFakeStr(10), DateOfBirth: shapeFakeStr(10), Nationality: shapeFakeStr(10), MaritalStatus: shapeFakeStr(15),
		NationalId: shapeFakeStr(15), PhotoUrl: shapeFakeStr(25),
		Email: shapeFakeStr(25), PhonePrimary: shapeFakeStr(15), PhoneSecondary: shapeFakeStr(15), AddressStreet: shapeFakeStr(20),
		AddressCity: shapeFakeStr(12), AddressProvince: shapeFakeStr(10), AddressPostalCode: shapeFakeStr(10), AddressCountry: shapeFakeStr(10),
		EmergencyContactName: shapeFakeStr(18), EmergencyContactPhone: shapeFakeStr(15),
		JobTitle: shapeFakeStr(18), DepartmentName: shapeFakeStr(15), TeamName: shapeFakeStr(15), ManagerName: shapeFakeStr(18),
		EmploymentStatus: shapeFakeStr(14), EmploymentType: shapeFakeStr(12), HireDate: shapeFakeStr(10), WorkLocation: shapeFakeStr(14),
		WorkMode: shapeFakeStr(10), BadgeNumber: shapeFakeStr(10),
		BaseSalary: shapeFakeInt64(30_000_000), Currency: shapeFakeStr(3), PayFrequency: shapeFakeStr(8), BonusTargetPct: shapeFakeFloat64(),
		BonusLastPaid: shapeFakeInt64(5_000_000), AllowanceTransport: shapeFakeInt64(1_000_000), AllowanceMeal: shapeFakeInt64(1_000_000),
		TaxId: shapeFakeStr(15), BankAccountLast4: shapeFakeStr(4), SalaryReviewDate: shapeFakeStr(10),
		PerformanceScore: shapeFakeFloat64(), PerformanceRating: shapeFakeStr(10), LastReviewDate: shapeFakeStr(10), NextReviewDate: shapeFakeStr(10),
		GoalsCompleted: shapeFakeInt32(20), GoalsTotal: shapeFakeInt32(20), PeerReviewScore: shapeFakeFloat64(), PromotionEligible: shapeFakeBool(),
		DisciplinaryActionsCount: shapeFakeInt32(3), CommendationsCount: shapeFakeInt32(5),
		AttendanceRate: shapeFakeFloat64(), LeaveBalanceDays: shapeFakeInt32(20), SickDaysUsed: shapeFakeInt32(10), VacationDaysUsed: shapeFakeInt32(15),
		UnpaidLeaveDays: shapeFakeInt32(5), LateArrivalsCount: shapeFakeInt32(10), RemoteDaysPerWeek: shapeFakeInt32(5),
		OvertimeHoursYtd: shapeFakeFloat64(), ShiftPattern: shapeFakeStr(10), Timezone: shapeFakeStr(12),
		PrimarySkill: shapeFakeStr(12), SecondarySkill: shapeFakeStr(12), CertificationCount: shapeFakeInt32(5), TrainingHoursYtd: shapeFakeFloat64(),
		LanguagePrimary: shapeFakeStr(10), EducationLevel: shapeFakeStr(14), UniversityName: shapeFakeStr(22), GraduationYear: shapeFakeInt32(2023),
		DegreeField: shapeFakeStr(15), Gpa: shapeFakeFloat64(),
		ProjectCountActive: shapeFakeInt32(5), ProjectCountCompleted: shapeFakeInt32(20), CurrentProjectName: shapeFakeStr(20),
		CurrentProjectRole: shapeFakeStr(14), UtilizationRate: shapeFakeFloat64(), BillableHoursYtd: shapeFakeFloat64(),
		ClientFacing: shapeFakeBool(), TravelRequired: shapeFakeBool(), OnCallRotation: shapeFakeBool(), SystemAccessLevel: shapeFakeStr(10),
		ProbationStatus: shapeFakeStr(10), VisaStatus: shapeFakeStr(10), VisaExpiryDate: shapeFakeStr(10), BackgroundCheckStatus: shapeFakeStr(12),
		NdaSigned: shapeFakeBool(), InsurancePlan: shapeFakeStr(12), DependentsCount: shapeFakeInt32(4), RetirementPlanEnrolled: shapeFakeBool(),
		RetirementContributionPct: shapeFakeFloat64(), ReferralSource: shapeFakeStr(14),
		YearsOfService: shapeFakeFloat64(), PreviousCompany: shapeFakeStr(20), PreviousJobTitle: shapeFakeStr(15), MentorName: shapeFakeStr(15),
		MenteeCount: shapeFakeInt32(3), CommitteeMemberships: shapeFakeInt32(2), VolunteerHours: shapeFakeFloat64(),
		WellnessProgramEnrolled: shapeFakeBool(), CommuteDistanceKm: shapeFakeFloat64(), EngagementScore: shapeFakeFloat64(),
	}
}

// ToShapeDepth0Response generates n flat, non-nested EmployeeFlat records.
func ToShapeDepth0Response(n int) *model.ShapeDepth0Response {
	employees := make([]*model.EmployeeFlat, n)
	for i := range employees {
		employees[i] = shapeFakeEmployeeFlat()
	}
	return &model.ShapeDepth0Response{Employees: employees}
}

// ---------- Depth 1, Wide ----------

func shapeFakeEmployeeThin() *model.EmployeeThin {
	return &model.EmployeeThin{
		EmployeeId: shapeFakeStr(10), FullName: shapeFakeStr(20), JobTitle: shapeFakeStr(15),
		HireDate: shapeFakeStr(10), PerformanceScore: shapeFakeFloat64(),
	}
}

func shapeFakeCompanyDepth1Wide() *model.CompanyDepth1Wide {
	employees := make([]*model.EmployeeThin, depth1WideEmployeesPerCompany)
	for i := range employees {
		employees[i] = shapeFakeEmployeeThin()
	}
	return &model.CompanyDepth1Wide{
		CompanyId: shapeFakeStr(10), CompanyName: shapeFakeStr(20), Industry: shapeFakeStr(15), FoundedYear: shapeFakeInt32(2015),
		HeadquartersCity: shapeFakeStr(12), EmployeeCountTotal: shapeFakeInt32(500), AnnualRevenue: shapeFakeInt64(5_000_000_000),
		CeoName: shapeFakeStr(18), Employees: employees,
	}
}

// ToShapeDepth1WideResponse generates n Company records, each holding a
// wide (40-element), single-level array of EmployeeThin.
func ToShapeDepth1WideResponse(n int) *model.ShapeDepth1WideResponse {
	companies := make([]*model.CompanyDepth1Wide, n)
	for i := range companies {
		companies[i] = shapeFakeCompanyDepth1Wide()
	}
	return &model.ShapeDepth1WideResponse{Companies: companies}
}

// ---------- Depth 3, Narrow ----------

func shapeFakeEmployeeRich() *model.EmployeeRich {
	return &model.EmployeeRich{
		EmployeeId: shapeFakeStr(10), FullName: shapeFakeStr(20), JobTitle: shapeFakeStr(15), Email: shapeFakeStr(25), Phone: shapeFakeStr(15),
		HireDate: shapeFakeStr(10), EmploymentStatus: shapeFakeStr(12), Currency: shapeFakeStr(3), WorkLocation: shapeFakeStr(14),
		ManagerName: shapeFakeStr(18), LastReviewDate: shapeFakeStr(10), BirthDate: shapeFakeStr(10), Gender: shapeFakeStr(10),
		Nationality: shapeFakeStr(10), EmergencyContactName: shapeFakeStr(15), EmergencyContactPhone: shapeFakeStr(15),
		BadgeNumber: shapeFakeStr(10), ContractEndDate: shapeFakeStr(10),
		BaseSalary: shapeFakeInt64(30_000_000), PerformanceScore: shapeFakeFloat64(), AttendanceRate: shapeFakeFloat64(),
		BonusTargetPct: shapeFakeFloat64(), OvertimeRate: shapeFakeFloat64(), TrainingHoursYtd: shapeFakeFloat64(),
		LeaveBalanceDays: shapeFakeInt32(20), ProjectCountActive: shapeFakeInt32(5),
	}
}

func shapeFakeTeamNarrow() *model.TeamNarrow {
	employees := make([]*model.EmployeeRich, depth3NarrowEmployeesPerTeam)
	for i := range employees {
		employees[i] = shapeFakeEmployeeRich()
	}
	return &model.TeamNarrow{
		TeamId: shapeFakeStr(10), TeamName: shapeFakeStr(18), TeamLead: shapeFakeStr(18), FocusArea: shapeFakeStr(15),
		FormedDate: shapeFakeStr(10), Employees: employees,
	}
}

func shapeFakeDepartmentNarrow() *model.DepartmentNarrow {
	teams := make([]*model.TeamNarrow, depth3NarrowTeamsPerDept)
	for i := range teams {
		teams[i] = shapeFakeTeamNarrow()
	}
	return &model.DepartmentNarrow{
		DepartmentId: shapeFakeStr(10), DepartmentName: shapeFakeStr(18), Headcount: shapeFakeInt32(42),
		BudgetAnnual: shapeFakeInt64(4_200_000_000), Location: shapeFakeStr(12), Teams: teams,
	}
}

func shapeFakeCompanyDepth3Narrow() *model.CompanyDepth3Narrow {
	departments := make([]*model.DepartmentNarrow, depth3NarrowDepartments)
	for i := range departments {
		departments[i] = shapeFakeDepartmentNarrow()
	}
	return &model.CompanyDepth3Narrow{
		CompanyId: shapeFakeStr(10), CompanyName: shapeFakeStr(20), Industry: shapeFakeStr(15), FoundedYear: shapeFakeInt32(2015),
		HeadquartersCity: shapeFakeStr(12), EmployeeCountTotal: shapeFakeInt32(500), AnnualRevenue: shapeFakeInt64(5_000_000_000),
		CeoName: shapeFakeStr(18), Departments: departments,
	}
}

// ToShapeDepth3NarrowResponse generates n Company records with the
// Company -> Department[] -> Team[] -> EmployeeRich[] chain (3 levels deep).
func ToShapeDepth3NarrowResponse(n int) *model.ShapeDepth3NarrowResponse {
	companies := make([]*model.CompanyDepth3Narrow, n)
	for i := range companies {
		companies[i] = shapeFakeCompanyDepth3Narrow()
	}
	return &model.ShapeDepth3NarrowResponse{Companies: companies}
}

// ---------- Depth 4, Wide ----------

func shapeFakeCertification() *model.Certification {
	return &model.Certification{
		CertificationId: shapeFakeStr(10), CertificationName: shapeFakeStr(30), IssuingBody: shapeFakeStr(20),
		IssuedDate: shapeFakeStr(10), Score: shapeFakeFloat64(),
	}
}

func shapeFakeEmployeeSlim() *model.EmployeeSlim {
	certs := make([]*model.Certification, depth4WideCertsPerEmployee)
	for i := range certs {
		certs[i] = shapeFakeCertification()
	}
	return &model.EmployeeSlim{
		EmployeeId: shapeFakeStr(10), FullName: shapeFakeStr(20), JobTitle: shapeFakeStr(15), HireDate: shapeFakeStr(10),
		TeamRole: shapeFakeStr(14), Certifications: certs,
	}
}

func shapeFakeTeamWide() *model.TeamWide {
	employees := make([]*model.EmployeeSlim, depth4WideEmployeesPerTeam)
	for i := range employees {
		employees[i] = shapeFakeEmployeeSlim()
	}
	return &model.TeamWide{
		TeamId: shapeFakeStr(10), TeamName: shapeFakeStr(18), TeamLead: shapeFakeStr(18), FocusArea: shapeFakeStr(15),
		FormedDate: shapeFakeStr(10), Employees: employees,
	}
}

func shapeFakeDepartmentWide() *model.DepartmentWide {
	teams := make([]*model.TeamWide, depth4WideTeamsPerDept)
	for i := range teams {
		teams[i] = shapeFakeTeamWide()
	}
	return &model.DepartmentWide{
		DepartmentId: shapeFakeStr(10), DepartmentName: shapeFakeStr(18), Headcount: shapeFakeInt32(42),
		BudgetAnnual: shapeFakeInt64(4_200_000_000), Location: shapeFakeStr(12), Teams: teams,
	}
}

func shapeFakeCompanyDepth4Wide() *model.CompanyDepth4Wide {
	departments := make([]*model.DepartmentWide, depth4WideDepartments)
	for i := range departments {
		departments[i] = shapeFakeDepartmentWide()
	}
	return &model.CompanyDepth4Wide{
		CompanyId: shapeFakeStr(10), CompanyName: shapeFakeStr(20), Industry: shapeFakeStr(15), FoundedYear: shapeFakeInt32(2015),
		HeadquartersCity: shapeFakeStr(12), EmployeeCountTotal: shapeFakeInt32(500), AnnualRevenue: shapeFakeInt64(5_000_000_000),
		CeoName: shapeFakeStr(18), Departments: departments,
	}
}

// ToShapeDepth4WideResponse generates n Company records with the
// Company -> Department[] -> Team[] -> EmployeeSlim[] -> Certification[]
// chain (4 levels deep), with the widest array at the innermost level.
func ToShapeDepth4WideResponse(n int) *model.ShapeDepth4WideResponse {
	companies := make([]*model.CompanyDepth4Wide, n)
	for i := range companies {
		companies[i] = shapeFakeCompanyDepth4Wide()
	}
	return &model.ShapeDepth4WideResponse{Companies: companies}
}
