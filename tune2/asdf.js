const [store, setStore] = useState([]);

const addHandler = (e) => {
  setStore([...store, e]);
};
