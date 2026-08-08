package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"

	"google.golang.org/protobuf/proto"

	"scaling-up-rest-vs-grpc/internal/data/model"
)

const (
	targetCompact = 100_000
	targetLarge   = 500_000

	depth1WideEmployeesPerCompany = 40

	depth3NarrowDepartments      = 3
	depth3NarrowTeamsPerDept     = 2
	depth3NarrowEmployeesPerTeam = 4

	depth4WideDepartments      = 2
	depth4WideTeamsPerDept     = 2
	depth4WideEmployeesPerTeam = 3
	depth4WideCertsPerEmployee = 15
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// fakeStr returns a random string of exactly n characters, used to keep
// per-field byte sizes consistent with the earlier analytical estimate.
func fakeStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func fakeInt32(max int32) int32 { return int32(rand.Intn(int(max)) + 1) }
func fakeInt64(max int64) int64 { return rand.Int63n(max) + 1 }
func fakeFloat64() float64      { return rand.Float64() * 100 }
func fakeBool() bool            { return rand.Intn(2) == 0 }

func main() {
	fmt.Printf("%-16s %12s %12s | %12s %12s\n", "Variasi", "N-compact", "byte", "N-large", "byte")

	nc, bc := calibrate(buildDepth0, targetCompact)
	nl, bl := calibrate(buildDepth0, targetLarge)
	fmt.Printf("%-16s %12d %12d | %12d %12d\n", "Depth0", nc, bc, nl, bl)

	nc, bc = calibrate(buildDepth1Wide, targetCompact)
	nl, bl = calibrate(buildDepth1Wide, targetLarge)
	fmt.Printf("%-16s %12d %12d | %12d %12d\n", "Depth1-Wide", nc, bc, nl, bl)

	nc, bc = calibrate(buildDepth3Narrow, targetCompact)
	nl, bl = calibrate(buildDepth3Narrow, targetLarge)
	fmt.Printf("%-16s %12d %12d | %12d %12d\n", "Depth3-Narrow", nc, bc, nl, bl)

	nc, bc = calibrate(buildDepth4Wide, targetCompact)
	nl, bl = calibrate(buildDepth4Wide, targetLarge)
	fmt.Printf("%-16s %12d %12d | %12d %12d\n", "Depth4-Wide", nc, bc, nl, bl)
}

// calibrate finds the N whose serialized size is closest to target,
// incrementing N by one document at a time until the total crosses target.
func calibrate(build func(n int) proto.Message, target int) (int, int) {
	bestN, bestSize, bestDiff := 1, 0, -1
	for n := 1; ; n++ {
		b, err := proto.Marshal(build(n))
		if err != nil {
			slog.Error("marshal failed", "error", err)
			os.Exit(1)
		}
		size := len(b)
		diff := size - target
		if diff < 0 {
			diff = -diff
		}
		if bestDiff == -1 || diff < bestDiff {
			bestDiff, bestN, bestSize = diff, n, size
		}
		if size > target {
			break
		}
	}
	return bestN, bestSize
}

// ---------- Depth 0 ----------

func fakeEmployeeFlat() *model.EmployeeFlat {
	return &model.EmployeeFlat{
		EmployeeId: fakeStr(10), FullName: fakeStr(18), FirstName: fakeStr(10), LastName: fakeStr(10),
		Gender: fakeStr(10), DateOfBirth: fakeStr(10), Nationality: fakeStr(10), MaritalStatus: fakeStr(15),
		NationalId: fakeStr(15), PhotoUrl: fakeStr(25),
		Email: fakeStr(25), PhonePrimary: fakeStr(15), PhoneSecondary: fakeStr(15), AddressStreet: fakeStr(20),
		AddressCity: fakeStr(12), AddressProvince: fakeStr(10), AddressPostalCode: fakeStr(10), AddressCountry: fakeStr(10),
		EmergencyContactName: fakeStr(18), EmergencyContactPhone: fakeStr(15),
		JobTitle: fakeStr(18), DepartmentName: fakeStr(15), TeamName: fakeStr(15), ManagerName: fakeStr(18),
		EmploymentStatus: fakeStr(14), EmploymentType: fakeStr(12), HireDate: fakeStr(10), WorkLocation: fakeStr(14),
		WorkMode: fakeStr(10), BadgeNumber: fakeStr(10),
		BaseSalary: fakeInt64(30_000_000), Currency: fakeStr(3), PayFrequency: fakeStr(8), BonusTargetPct: fakeFloat64(),
		BonusLastPaid: fakeInt64(5_000_000), AllowanceTransport: fakeInt64(1_000_000), AllowanceMeal: fakeInt64(1_000_000),
		TaxId: fakeStr(15), BankAccountLast4: fakeStr(4), SalaryReviewDate: fakeStr(10),
		PerformanceScore: fakeFloat64(), PerformanceRating: fakeStr(10), LastReviewDate: fakeStr(10), NextReviewDate: fakeStr(10),
		GoalsCompleted: fakeInt32(20), GoalsTotal: fakeInt32(20), PeerReviewScore: fakeFloat64(), PromotionEligible: fakeBool(),
		DisciplinaryActionsCount: fakeInt32(3), CommendationsCount: fakeInt32(5),
		AttendanceRate: fakeFloat64(), LeaveBalanceDays: fakeInt32(20), SickDaysUsed: fakeInt32(10), VacationDaysUsed: fakeInt32(15),
		UnpaidLeaveDays: fakeInt32(5), LateArrivalsCount: fakeInt32(10), RemoteDaysPerWeek: fakeInt32(5),
		OvertimeHoursYtd: fakeFloat64(), ShiftPattern: fakeStr(10), Timezone: fakeStr(12),
		PrimarySkill: fakeStr(12), SecondarySkill: fakeStr(12), CertificationCount: fakeInt32(5), TrainingHoursYtd: fakeFloat64(),
		LanguagePrimary: fakeStr(10), EducationLevel: fakeStr(14), UniversityName: fakeStr(22), GraduationYear: fakeInt32(2023),
		DegreeField: fakeStr(15), Gpa: fakeFloat64(),
		ProjectCountActive: fakeInt32(5), ProjectCountCompleted: fakeInt32(20), CurrentProjectName: fakeStr(20),
		CurrentProjectRole: fakeStr(14), UtilizationRate: fakeFloat64(), BillableHoursYtd: fakeFloat64(),
		ClientFacing: fakeBool(), TravelRequired: fakeBool(), OnCallRotation: fakeBool(), SystemAccessLevel: fakeStr(10),
		ProbationStatus: fakeStr(10), VisaStatus: fakeStr(10), VisaExpiryDate: fakeStr(10), BackgroundCheckStatus: fakeStr(12),
		NdaSigned: fakeBool(), InsurancePlan: fakeStr(12), DependentsCount: fakeInt32(4), RetirementPlanEnrolled: fakeBool(),
		RetirementContributionPct: fakeFloat64(), ReferralSource: fakeStr(14),
		YearsOfService: fakeFloat64(), PreviousCompany: fakeStr(20), PreviousJobTitle: fakeStr(15), MentorName: fakeStr(15),
		MenteeCount: fakeInt32(3), CommitteeMemberships: fakeInt32(2), VolunteerHours: fakeFloat64(),
		WellnessProgramEnrolled: fakeBool(), CommuteDistanceKm: fakeFloat64(), EngagementScore: fakeFloat64(),
	}
}

func buildDepth0(n int) proto.Message {
	employees := make([]*model.EmployeeFlat, n)
	for i := range employees {
		employees[i] = fakeEmployeeFlat()
	}
	return &model.ShapeDepth0Response{Employees: employees}
}

// ---------- Depth 1, Wide ----------

func fakeEmployeeThin() *model.EmployeeThin {
	return &model.EmployeeThin{
		EmployeeId: fakeStr(10), FullName: fakeStr(20), JobTitle: fakeStr(15),
		HireDate: fakeStr(10), PerformanceScore: fakeFloat64(),
	}
}

func fakeCompanyDepth1Wide() *model.CompanyDepth1Wide {
	employees := make([]*model.EmployeeThin, depth1WideEmployeesPerCompany)
	for i := range employees {
		employees[i] = fakeEmployeeThin()
	}
	return &model.CompanyDepth1Wide{
		CompanyId: fakeStr(10), CompanyName: fakeStr(20), Industry: fakeStr(15), FoundedYear: fakeInt32(2015),
		HeadquartersCity: fakeStr(12), EmployeeCountTotal: fakeInt32(500), AnnualRevenue: fakeInt64(5_000_000_000),
		CeoName: fakeStr(18), Employees: employees,
	}
}

func buildDepth1Wide(n int) proto.Message {
	companies := make([]*model.CompanyDepth1Wide, n)
	for i := range companies {
		companies[i] = fakeCompanyDepth1Wide()
	}
	return &model.ShapeDepth1WideResponse{Companies: companies}
}

// ---------- Depth 3, Narrow ----------

func fakeEmployeeRich() *model.EmployeeRich {
	return &model.EmployeeRich{
		EmployeeId: fakeStr(10), FullName: fakeStr(20), JobTitle: fakeStr(15), Email: fakeStr(25), Phone: fakeStr(15),
		HireDate: fakeStr(10), EmploymentStatus: fakeStr(12), Currency: fakeStr(3), WorkLocation: fakeStr(14),
		ManagerName: fakeStr(18), LastReviewDate: fakeStr(10), BirthDate: fakeStr(10), Gender: fakeStr(10),
		Nationality: fakeStr(10), EmergencyContactName: fakeStr(15), EmergencyContactPhone: fakeStr(15),
		BadgeNumber: fakeStr(10), ContractEndDate: fakeStr(10),
		BaseSalary: fakeInt64(30_000_000), PerformanceScore: fakeFloat64(), AttendanceRate: fakeFloat64(),
		BonusTargetPct: fakeFloat64(), OvertimeRate: fakeFloat64(), TrainingHoursYtd: fakeFloat64(),
		LeaveBalanceDays: fakeInt32(20), ProjectCountActive: fakeInt32(5),
	}
}

func fakeTeamNarrow() *model.TeamNarrow {
	employees := make([]*model.EmployeeRich, depth3NarrowEmployeesPerTeam)
	for i := range employees {
		employees[i] = fakeEmployeeRich()
	}
	return &model.TeamNarrow{
		TeamId: fakeStr(10), TeamName: fakeStr(18), TeamLead: fakeStr(18), FocusArea: fakeStr(15),
		FormedDate: fakeStr(10), Employees: employees,
	}
}

func fakeDepartmentNarrow() *model.DepartmentNarrow {
	teams := make([]*model.TeamNarrow, depth3NarrowTeamsPerDept)
	for i := range teams {
		teams[i] = fakeTeamNarrow()
	}
	return &model.DepartmentNarrow{
		DepartmentId: fakeStr(10), DepartmentName: fakeStr(18), Headcount: fakeInt32(42),
		BudgetAnnual: fakeInt64(4_200_000_000), Location: fakeStr(12), Teams: teams,
	}
}

func fakeCompanyDepth3Narrow() *model.CompanyDepth3Narrow {
	departments := make([]*model.DepartmentNarrow, depth3NarrowDepartments)
	for i := range departments {
		departments[i] = fakeDepartmentNarrow()
	}
	return &model.CompanyDepth3Narrow{
		CompanyId: fakeStr(10), CompanyName: fakeStr(20), Industry: fakeStr(15), FoundedYear: fakeInt32(2015),
		HeadquartersCity: fakeStr(12), EmployeeCountTotal: fakeInt32(500), AnnualRevenue: fakeInt64(5_000_000_000),
		CeoName: fakeStr(18), Departments: departments,
	}
}

func buildDepth3Narrow(n int) proto.Message {
	companies := make([]*model.CompanyDepth3Narrow, n)
	for i := range companies {
		companies[i] = fakeCompanyDepth3Narrow()
	}
	return &model.ShapeDepth3NarrowResponse{Companies: companies}
}

// ---------- Depth 4, Wide ----------

func fakeCertification() *model.Certification {
	return &model.Certification{
		CertificationId: fakeStr(10), CertificationName: fakeStr(30), IssuingBody: fakeStr(20),
		IssuedDate: fakeStr(10), Score: fakeFloat64(),
	}
}

func fakeEmployeeSlim() *model.EmployeeSlim {
	certs := make([]*model.Certification, depth4WideCertsPerEmployee)
	for i := range certs {
		certs[i] = fakeCertification()
	}
	return &model.EmployeeSlim{
		EmployeeId: fakeStr(10), FullName: fakeStr(20), JobTitle: fakeStr(15), HireDate: fakeStr(10),
		TeamRole: fakeStr(14), Certifications: certs,
	}
}

func fakeTeamWide() *model.TeamWide {
	employees := make([]*model.EmployeeSlim, depth4WideEmployeesPerTeam)
	for i := range employees {
		employees[i] = fakeEmployeeSlim()
	}
	return &model.TeamWide{
		TeamId: fakeStr(10), TeamName: fakeStr(18), TeamLead: fakeStr(18), FocusArea: fakeStr(15),
		FormedDate: fakeStr(10), Employees: employees,
	}
}

func fakeDepartmentWide() *model.DepartmentWide {
	teams := make([]*model.TeamWide, depth4WideTeamsPerDept)
	for i := range teams {
		teams[i] = fakeTeamWide()
	}
	return &model.DepartmentWide{
		DepartmentId: fakeStr(10), DepartmentName: fakeStr(18), Headcount: fakeInt32(42),
		BudgetAnnual: fakeInt64(4_200_000_000), Location: fakeStr(12), Teams: teams,
	}
}

func fakeCompanyDepth4Wide() *model.CompanyDepth4Wide {
	departments := make([]*model.DepartmentWide, depth4WideDepartments)
	for i := range departments {
		departments[i] = fakeDepartmentWide()
	}
	return &model.CompanyDepth4Wide{
		CompanyId: fakeStr(10), CompanyName: fakeStr(20), Industry: fakeStr(15), FoundedYear: fakeInt32(2015),
		HeadquartersCity: fakeStr(12), EmployeeCountTotal: fakeInt32(500), AnnualRevenue: fakeInt64(5_000_000_000),
		CeoName: fakeStr(18), Departments: departments,
	}
}

func buildDepth4Wide(n int) proto.Message {
	companies := make([]*model.CompanyDepth4Wide, n)
	for i := range companies {
		companies[i] = fakeCompanyDepth4Wide()
	}
	return &model.ShapeDepth4WideResponse{Companies: companies}
}
