package evedb

type Icon struct {
	IconID          int
	IconFile        string
	IconDescription string
}

// A Race is a race in EVE.
type Race struct {
	ID               int
	Name             string
	Description      string
	ShortDescription string

	Icon
}

type Ancestry struct {
	ID               int
	Name             string
	Description      string
	BloodlineID      int
	Perception       int
	Willpower        int
	Charisma         int
	Memory           int
	Intelligence     int
	ShortDescription string

	Icon
}

type Bloodline struct {
	ID                     int
	Name                   string
	RaceID                 int
	Description            string
	MaleDescription        string
	FemaleDescription      string
	ShipTypeID             int
	CorporationID          int
	Perception             int
	Willpower              int
	Charisma               int
	Memory                 int
	Intelligence           int
	ShortDescription       string
	ShortMaleDescription   string
	ShortFemaleDescription string

	Icon
}

func (e *EveDB) GetRace(id int) (*Race, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	row := c.QueryRow(
		`SELECT
			  race."raceID"
			, race."raceName"
			, COALESCE(race."description", '')
			, COALESCE(icon."iconFile", '')
			, COALESCE(race."shortDescription", '')
			FROM evesde."chrRaces" race
			LEFT JOIN evesde."eveIcons" icon ON race."iconID" = icon."iconID"
			WHERE race."raceID" = $1`, id)
	r := Race{Icon: Icon{}}
	err = row.Scan(&r.ID, &r.Name, &r.Description, &r.IconFile, &r.ShortDescription)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// GetRaces fetches all Races from the database.
func (e *EveDB) GetRaces() ([]*Race, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	rs, err := c.Query(
		`SELECT
			  race."raceID"
			, race."raceName"
			, COALESCE(race."description", '')
			, COALESCE(icon."iconFile", '')
			, COALESCE(race."shortDescription", '')
			FROM evesde."chrRaces" race
			LEFT JOIN evesde."eveIcons" icon ON race."iconID" = icon."iconID"`)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	res := []*Race{}
	for rs.Next() {
		r := Race{Icon: Icon{}}
		err := rs.Scan(&r.ID, &r.Name, &r.Description, &r.IconFile, &r.ShortDescription)
		if err != nil {
			return nil, err
		}
		res = append(res, &r)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (e *EveDB) GetAncestry(id int) (*Ancestry, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	row := c.QueryRow(
		`SELECT
			  ancestry."ancestryID"
			, ancestry."ancestryName"
			, COALESCE(ancestry."description", '')
			, COALESCE(icon."iconFile", '')
			, COALESCE(ancestry."shortDescription", '')
			FROM evesde."chrAncestries" ancestry
			LEFT JOIN evesde."eveIcons" icon ON ancestry."iconID" = icon."iconID"
			WHERE ancestry."ancestryID" = $1`, id)
	a := Ancestry{Icon: Icon{}}
	err = row.Scan(&a.ID, &a.Name, &a.Description, &a.IconFile, &a.ShortDescription)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (e *EveDB) GetBloodline(id int) (*Bloodline, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	row := c.QueryRow(
		`SELECT
			  bloodline."bloodlineID"
			, bloodline."bloodlineName"
			, COALESCE(bloodline."description", '')
			, COALESCE(icon."iconFile", '')
			, COALESCE(bloodline."shortDescription", '')
			FROM evesde."chrBloodlines" bloodline
			LEFT JOIN evesde."eveIcons" icon ON bloodline."iconID" = icon."iconID"
			WHERE bloodline."bloodlineID" = $1`, id)
	b := Bloodline{Icon: Icon{}}
	err = row.Scan(&b.ID, &b.Name, &b.Description, &b.IconFile, &b.ShortDescription)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
