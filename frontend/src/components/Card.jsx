export default function Card({ children, className = '', ...props }) {
  return (
    <div
      className={`rounded-[28px] border border-[#948979]/20 bg-[#393E46] shadow-2xl shadow-black/20 ${className}`.trim()}
      {...props}
    >
      {children}
    </div>
  );
}
