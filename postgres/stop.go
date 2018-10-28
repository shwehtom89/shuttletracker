package postgres

import (
	"database/sql"

	"github.com/wtg/shuttletracker"
)

// StopService is an implementation of shuttletracker.StopService.
type StopService struct {
	db *sql.DB
}

func (ss *StopService) initializeSchema(db *sql.DB) error {
	ss.db = db
	schema := `
CREATE TABLE IF NOT EXISTS stops (
	id serial PRIMARY KEY,
	name text,
	description text,
	latitude double precision NOT NULL,
	longitude double precision NOT NULL,
	created timestamp with time zone NOT NULL DEFAULT now(),
	updated timestamp with time zone NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS schedules (
	id serial PRIMARY KEY,
	name text NOT NULL,
	weekend boolean NOT NULL,
	west boolean NOT NULL
);
CREATE TABLE IF NOT EXISTS schedule_times (
	id serial PRIMARY KEY,
	schedule_id integer REFERENCES schedules NOT NULL,
	stop_id integer REFERENCES stops NOT NULL,
	time integer NOT NULL
);`
	_, err := ss.db.Exec(schema)
	return err
}

// CreateStop creates a Stop.
func (ss *StopService) CreateStop(stop *shuttletracker.Stop) error {
	statement := "INSERT INTO stops (name, description, latitude, longitude) VALUES" +
		" ($1, $2, $3, $4) RETURNING id, created, updated;"
	row := ss.db.QueryRow(statement, stop.Name, stop.Description, stop.Latitude, stop.Longitude)
	return row.Scan(&stop.ID, &stop.Created, &stop.Updated)
}

// Stops returns all Stops.
func (ss *StopService) Stops() ([]*shuttletracker.Stop, error) {
	stops := []*shuttletracker.Stop{}
	query := "SELECT s.id, s.name, s.created, s.updated, s.description, s.latitude, s.longitude" +
		" FROM stops s;"
	rows, err := ss.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &shuttletracker.Stop{}
		err := rows.Scan(&s.ID, &s.Name, &s.Created, &s.Updated, &s.Description, &s.Latitude, &s.Longitude)
		if err != nil {
			return nil, err
		}
		stops = append(stops, s)
	}
	return stops, nil
}

// DeleteStop deletes a Stop.
func (ss *StopService) DeleteStop(id int64) error {
	statement := "DELETE FROM stops WHERE id = $1;"
	result, err := ss.db.Exec(statement, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return shuttletracker.ErrStopNotFound
	}

	return nil
}

// ScheduleStops returns all ScheduleStops associated with one schedule.
func (ss *StopService) ScheduleStops(id int64) ([]*shuttletracker.ScheduleStop, error) {
	return []*shuttletracker.ScheduleStop{}, nil
}

// use to insert json from parsed excel file
func (ss *StopService) InsertScheduleStops(stops []interface{}) {
	// unmarshal data into map[string]interface{}
	// m := []interface{}{}
	// err := json.Unmarshal(data, &m)

	// for i, entry := range stops {
		//     find stop name in stops and set stop_id foreign key
		//     set schedule name
		//     set time
		//     make sure it gets inserted into the database
	// }
}
