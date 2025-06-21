import "./test.css";

type WallProps = {
  horizontal: boolean;
};

export const Wall = ({ horizontal = false }: WallProps) => {
  return (
    <div
      style={{
        position: "relative",
        width: "32px",
        height: "32px",
        animation: "Bounce",
      }}
    >
      <img src="floor.svg" alt="floor" />

      <div
        style={{
          position: "absolute",
          transform: horizontal ? "rotate(90)" : "",
        }}
      >
        <img src="wall.svg" alt="wall" />
      </div>
    </div>
  );
};

export const Floor = () => {
  return <img src="floor.svg" alt="floor" />;
};
