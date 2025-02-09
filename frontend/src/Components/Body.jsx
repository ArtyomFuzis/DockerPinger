import "./Body.css"
import React, {useEffect} from "react";
function formatDate(date) {
    let new_date = new Date(date)
    if (new_date.getTime() <= 0){
        return "Нет данных"
    }
    return new_date.toLocaleString()
}
function Body() {
    const [state, setState] = React.useState({})
    useEffect(() => {
        const interval = setInterval(() => {
            fetch('/api/info', {
                method: 'Get',
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Сетевая ошибка')
                    }
                    return response.json()
                })
                .then(data => {
                    setState((state) =>
                        ({...state, services: data})
                    )
                })
                .catch(error => console.error('Ошибка:', error))
        }, 1000)
        return () => clearInterval(interval);
    }, [state]);
    let table
    if (state.services) {
        table = state.services.map((el) => <tr key={el.Address}>
            {el.State ?
                <td>
                    <div className="status-ok status">OK</div>
                </td>
                :
                <td>
                    <div className="status-fail status">FAIL</div>
                </td>
            }
            <td>{el.Address}</td>
            <td>{formatDate(el.LastPing.Date)}</td>
            <td>{formatDate(el.LastSuccess.Date) }</td>
        </tr>)
    }
    return (
        <div className="App-body">
            <div className="control-pane">
                <input className="control-pane-element control-pane-input" type="text"
                       value={state.txt ? state.txt : ""} onChange={(e) => setState({...state, txt: e.target.value})}/>
                <input className="control-pane-element control-pane-btn control-pane-btn-add" type="submit"
                       value="Добавить" onClick={
                    () =>
                        fetch('/api/add?address=' + state.txt, {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/x-www-form-urlencoded',
                            }
                        })
                            .then(response => {
                                if (!response.ok) {
                                    throw new Error('Сетевая ошибка')
                                }
                            })
                            .catch(error => console.error('Ошибка:', error))
                }/>
                <input className="control-pane-element control-pane-btn control-pane-btn-delete" type="submit"
                       value="Удалить" onClick={
                    () =>
                        fetch('/api/delete?address=' + state.txt, {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/x-www-form-urlencoded',
                            }
                        })
                            .then(response => {
                                if (!response.ok) {
                                    throw new Error('Сетевая ошибка')
                                }
                            })
                            .catch(error => console.error('Ошибка:', error))
                }/>
            </div>
            <div className="table-pane">
                <table className="table">
                    <thead>
                    <tr>
                        <th>Состояние</th>
                        <th>Адрес</th>
                        <th>Последний пинг</th>
                        <th>Последний удачный пинг</th>
                    </tr>
                    </thead>
                    <tbody>
                    {table}
                    </tbody>
                </table>
            </div>
        </div>
    )
}

export default Body