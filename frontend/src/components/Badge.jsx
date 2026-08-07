export default function Badge({ children, className = '', ...props }) {
  return (
    <span
      className={`inline-flex rounded-full border border-[#948979]/30 bg-[#948979]/10 px-3 py-1 text-sm font-medium text-[#DFD0B8] ${className}`.trim()}
      {...props}
    >
      {children}
    </span>
  );
}
